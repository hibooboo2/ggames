package pollen

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/logger"
	"github.com/hibooboo2/ggames/allplay/pollen/cards"
	"github.com/hibooboo2/ggames/allplay/pollen/colors"
	"github.com/hibooboo2/ggames/allplay/pollen/position"
	"github.com/hibooboo2/ggames/allplay/pollen/token"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Game struct {
	id                 uuid.UUID
	Name               string
	Owner              string
	Players            []*Player
	InvitedUsers       map[string]struct{}
	PlayerUsernames    map[string]struct{}
	ActivePlayerCursor int
	TokenBag           *token.TokenBag
	Board              *Board
	events             chan struct{}
	done               chan struct{}
	started            bool
	readyChan          chan struct{}
}

func NewGame(id uuid.UUID, username string, gameName string) *Game {
	g := &Game{
		id:              id,
		Owner:           username,
		Name:            gameName,
		TokenBag:        token.NewTokenBag(true),
		events:          make(chan struct{}),
		done:            make(chan struct{}),
		readyChan:       make(chan struct{}),
		PlayerUsernames: map[string]struct{}{},
		InvitedUsers:    map[string]struct{}{},
	}
	g.AddPlayer(username)
	return g
}

func (g *Game) AllPlayersOutOfCards() bool {
	out := 0
	for _, player := range g.Players {
		if player.OutOfCards() {
			out++
		}
	}
	return out == len(g.Players)
}

func (g *Game) Start() error {
	select {
	case <-g.readyChan:
		return nil
	case <-time.After(time.Second):
	}
	switch len(g.PlayerUsernames) {
	case 2, 3, 4:
	default:
		return fmt.Errorf("invalid number of players")
	}

	usernames := []string{}
	for username := range g.PlayerUsernames {
		usernames = append(usernames, username)
	}
	sort.Strings(usernames)

	c := colors.Purple
	for _, user := range usernames {
		g.Players = append(g.Players, NewPlayer(user, len(g.PlayerUsernames), c))
		c = 1 << c
	}

	for _, p := range g.Players {
		logger.Games(p.Username, p.Color)
	}
	tk := g.TokenBag.TakeNextToken()
	logger.AtLevel(logger.LBoard|logger.LGames|logger.LToken, fmt.Sprintf("Token %v retrived from bag", tk))
	g.Board = NewBoard(tk, g)
	go func() {
		for range g.events {
			for _, p := range g.Players {
				p := p
				go func() {
					select {
					case p.events <- struct{}{}:
					case <-time.After(time.Millisecond * 100):
					case <-g.done:
						return
					}
				}()
			}
		}
	}()
	g.started = true
	close(g.readyChan)
	return nil
}

func (g *Game) Started() bool {
	return g.started
}

func (g *Game) AddPlayer(username string) error {
	if len(g.PlayerUsernames) > 3 {
		return errors.New("game already has 4 players")
	}

	logger.Gamesf("Adding user to game: %s %s", g.id, username)
	g.PlayerUsernames[username] = struct{}{}
	delete(g.InvitedUsers, username)
	return nil
}

func (g *Game) InvitePlayer(username string) {
	g.InvitedUsers[username] = struct{}{}
}

func (g *Game) HasPlayer(username string) bool {
	for uname := range g.PlayerUsernames {
		logger.Gamesf("Checking if user is in game: %s %q", g.id, uname)
		if uname == username {
			return true
		}
	}
	return false
}
func (g *Game) HasInvitedUser(username string) bool {
	_, ok := g.InvitedUsers[username]
	return ok
}

func (g *Game) IsOwner(username string) bool {
	return g.Owner == username
}

func (g *Game) End() {
	close(g.done)
}

func (g *Game) ToggleHints(username string) {
	for _, p := range g.Players {
		if p.Username == username {
			p.HintsOn = !p.HintsOn
		}
	}
	logger.Gamesf("Hints toggled for user: %s", username)
	g.events <- struct{}{}
}

func (g *Game) GetID() uuid.UUID {
	return g.id
}

func (g *Game) GetHand(username string) []cards.GardenCard {
	logger.Gamesln("Getting hand for", username)
	for _, player := range g.Players {
		if player.Username == username {
			return player.Hand
		}
	}
	return nil
}

func (g *Game) NextPlayer() error {
	if g.Board.GameOver() {
		return errors.New("game is over")
	}

	if err := g.Board.MustPlayToken(); err != nil {
		return fmt.Errorf("current player must place a token: %w", err)
	}

	g.activePlayer().CardNotPlayed()
	g.ActivePlayerCursor += 1
	g.events <- struct{}{}
	return nil
}

func (g *Game) MustPlayTokens() []position.Position {
	if g.Board.MustPlayToken() == nil {
		return nil
	}
	return g.Board.GetTokensMustPlay()
}

func (g *Game) PlayCard(username string, card uuid.UUID, position position.Position) error {
	defer func() { g.events <- struct{}{} }()
	p := g.activePlayer()
	if p.Username != username {
		return fmt.Errorf("it is not %q turn as it is %q turn", username, p.Username)
	}
	cardToPlay, err := p.GetCard(card)
	if err != nil {
		return fmt.Errorf("error moving card: %w", err)
	}

	if p.CardPlayed() {
		return fmt.Errorf("player has already gone")
	}

	if err := g.Board.PlayCard(position, cardToPlay); err != nil {
		return fmt.Errorf("error playing card: %w", err)
	}

	_, err = p.PlayCard(cardToPlay.ID)
	if err != nil {
		return fmt.Errorf("error playing card: %w", err)
	}

	return nil
}

func (g *Game) activePlayer() *Player {
	return g.Players[g.ActivePlayerCursor%len(g.Players)]
}

func (g *Game) GetNextTokenID() *uuid.UUID {
	tokens := g.TokenBag.GetTokens(1)

	if tokens == nil {
		// This is endgame condition need to handle this
		return nil
	}

	return &tokens[0].ID
}

func (g *Game) PlayToken(username string, tokenID uuid.UUID, position position.Position) error {
	if g.activePlayer().Username != username {
		return fmt.Errorf("it is not %q turn", username)
	}

	if !g.TokenBag.HasToken(tokenID) {
		return fmt.Errorf("token not found. Was it already played?")
	}

	tk := g.TokenBag.GetToken(tokenID)
	if tk == nil {
		return fmt.Errorf("token was nil when retrieved")
	}

	if err := g.Board.PlayToken(position, tk); err != nil {
		return fmt.Errorf("failed to play token: %w", err)
	}

	g.events <- struct{}{}
	return nil
}

type FlusherWriter interface {
	io.Writer
	http.Flusher
}

func (g *Game) Render(c context.Context, w FlusherWriter, username string) error {
	var playerToRenderFor *Player
	for _, player := range g.Players {
		if player.Username == username {
			playerToRenderFor = player
		}
	}

	if playerToRenderFor == nil {
		_, isjoined := g.PlayerUsernames[username]
		_, invited := g.InvitedUsers[username]
		if !invited && !isjoined {
			return fmt.Errorf("player %q not found", username)
		}

		timer := time.After(time.Second * 10)

		fmt.Fprintf(w, "data: game_not_started\n\n")
		w.Flush()

		for playerToRenderFor == nil {
			select {
			case <-c.Done():
				logger.Gamesf("User %q disconnected", username)
				return nil
			case <-timer:
				return fmt.Errorf("player not found or game not started")
			case <-g.readyChan:
				for _, player := range g.Players {
					if player.Username == username {
						playerToRenderFor = player
					}
				}
				if playerToRenderFor == nil {
					return fmt.Errorf("player not a member of the game")
				}
			}
		}
	}

	fmt.Fprintf(w, "data: waiting\n\n")
	w.Flush()

	<-g.readyChan

	logger.Gamesln("Starting render for", playerToRenderFor.Username)

	err := g.Board.Render(w, playerToRenderFor, g)
	if err != nil {
		logger.Gamesf("Failed to render for %s: %v", playerToRenderFor.Username, err)
		fmt.Fprintf(w, "data: error %q\n\n", err.Error())
	}

	w.Flush()

	playerToRenderFor.ToggleConnection()
	defer playerToRenderFor.ToggleConnection()
	g.events <- struct{}{}
	for {
		select {
		case <-c.Done():
			logger.Gamesf("Player %q just disconnected", playerToRenderFor.Username)
			return nil
		case <-playerToRenderFor.events:
			err := g.Board.Render(w, playerToRenderFor, g)
			if err != nil {
				logger.Gamesf("Failed to render for %s: %v", playerToRenderFor.Username, err)
				fmt.Fprintf(w, "data: error %q\n\n", err.Error())
			}

			w.Flush()
		case <-g.done:
			w.Flush()
			return nil
		}
	}
}
