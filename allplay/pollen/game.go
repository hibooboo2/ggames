package pollen

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/logger"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Game struct {
	id                 uuid.UUID
	name               string
	owner              string
	players            []*Player
	invitedUsers       map[string]struct{}
	playerUsernames    map[string]struct{}
	activePlayerCursor int
	tokenBag           *TokenBag
	board              *Board
	events             chan struct{}
	done               chan struct{}
	started            bool
	readyChan          chan struct{}
	AutoToken          bool
}

func NewGame(id uuid.UUID, username string, gameName string, i int) *Game {
	g := &Game{
		id:              id,
		owner:           username,
		name:            fmt.Sprintf("%q %d", gameName, i),
		tokenBag:        NewTokenBag(),
		events:          make(chan struct{}),
		done:            make(chan struct{}),
		readyChan:       make(chan struct{}),
		playerUsernames: map[string]struct{}{},
		invitedUsers:    map[string]struct{}{},
	}
	g.AddPlayer(username)
	return g
}

func (g *Game) Name() string {
	return g.name
}

func (g *Game) Start() error {
	switch len(g.playerUsernames) {
	case 2, 3, 4:
	default:
		return fmt.Errorf("invalid number of players")
	}

	usernames := []string{}
	for username := range g.playerUsernames {
		usernames = append(usernames, username)
	}
	sort.Strings(usernames)

	for i, user := range usernames {
		g.players = append(g.players, NewPlayer(user, len(g.playerUsernames), Color(i)))
	}

	for _, p := range g.players {
		logger.Games(p.Username, p.Color)
	}
	g.board = NewBoard(g.tokenBag.GetToken())
	go func() {
		for range g.events {
			for _, p := range g.players {
				select {
				case p.Events <- struct{}{}:
				case <-time.After(time.Millisecond * 50):
				}
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
	if len(g.playerUsernames) > 3 {
		return errors.New("game already has 4 players")
	}

	logger.Gamesf("Adding user to game: %s %s", g.id, username)
	g.playerUsernames[username] = struct{}{}
	delete(g.invitedUsers, username)
	return nil
}

func (g *Game) InvitePlayer(username string) {
	g.invitedUsers[username] = struct{}{}
}

func (g *Game) HasPlayer(username string) bool {
	for uname := range g.playerUsernames {
		logger.Gamesf("Checking if user is in game: %s %q", g.id, uname)
		if uname == username {
			return true
		}
	}
	return false
}
func (g *Game) HasInvitedUser(username string) bool {
	_, ok := g.invitedUsers[username]
	return ok
}

func (g *Game) IsOwner(username string) bool {
	return g.owner == username
}

func (g *Game) End() {
	close(g.done)
}

func (g *Game) ToggleHints(username string) {
	for _, p := range g.players {
		if p.Username == username {
			p.HintsOn = !p.HintsOn
		}
	}
	g.events <- struct{}{}
}

func (g *Game) GetID() uuid.UUID {
	return g.id
}

func (g *Game) GetHand(username string) []GardenCard {
	logger.Gamesln("Getting hand for", username)
	for _, player := range g.players {
		if player.Username == username {
			return player.Hand
		}
	}
	return nil
}

func (g *Game) NextPlayer() error {
	if err := g.board.MustPlayToken(); err != nil {
		return fmt.Errorf("current player must place a token: %w", err)
	}
	g.activePlayerCursor += 1
	g.events <- struct{}{}
	return nil
}

func (g *Game) MustPlayTokens() []Position {
	if g.board.MustPlayToken() == nil {
		return nil
	}
	return g.board.GetTokensMustPlay()
}

func (g *Game) PlayCard(username string, card uuid.UUID, position Position) error {
	defer func() { g.events <- struct{}{} }()
	p := g.activePlayer()
	if p.Username != username {
		return fmt.Errorf("it is not %q turn as it is %q turn", username, p.Username)
	}
	cardToPlay, err := p.GetCard(card)
	if err != nil {
		return fmt.Errorf("error moving card: %w", err)
	}

	if err := g.board.PlayCard(position, cardToPlay); err != nil {
		return fmt.Errorf("error playing card: %w", err)
	}

	_, err = p.PlayCard(cardToPlay.ID)
	if err != nil {
		return fmt.Errorf("error playing card: %w", err)
	}

	//XXX make rest api for this and ui calls / clicks
	if g.AutoToken {
		tokenPositions := g.board.GetTokensMustPlay()
		for _, tokenPosition := range tokenPositions {
			tk := g.tokenBag.GetToken()
			err = g.PlayToken(username, tk, tokenPosition)
			if err != nil {
				return fmt.Errorf("failed to play token: %w", err)
			}
		}
	}

	return nil
}

func (g *Game) activePlayer() *Player {
	return g.players[g.activePlayerCursor%len(g.players)]
}

func (g *Game) GetNextToken() *PollinatorToken {
	token := g.tokenBag.GetToken()

	if token == nil {
		// This is endgame condition need to handle this
		return nil
	}

	return token
}

func (g *Game) PlayToken(username string, token *PollinatorToken, position Position) error {
	if g.activePlayer().Username != username {
		return fmt.Errorf("it is not %q turn", username)
	}

	if err := g.board.PlayToken(position, token); err != nil {
		return fmt.Errorf("failed to play token: %w", err)
	}

	g.events <- struct{}{}
	return nil
}

func (g *Game) Render(w io.Writer, username string) error {
	var playerToRenderFor *Player
	for _, player := range g.players {
		if player.Username == username {
			playerToRenderFor = player
		}
	}
	flusher := w.(http.Flusher)

	if playerToRenderFor == nil {
		_, isjoined := g.playerUsernames[username]
		_, invited := g.invitedUsers[username]
		if !invited && !isjoined {
			return fmt.Errorf("player %q not found", username)
		}

		timer := time.After(time.Minute * 5)

		fmt.Fprintf(w, "data: waiting\n\n")
		flusher.Flush()

		for playerToRenderFor == nil {
			select {
			case <-timer:
				return fmt.Errorf("player not found or game not started")
			case <-g.readyChan:
				for _, player := range g.players {
					if player.Username == username {
						playerToRenderFor = player
					}
				}
			}
		}
	}

	fmt.Fprintf(w, "data: waiting\n\n")
	flusher.Flush()

	<-g.readyChan

	logger.Gamesln("Starting render for", playerToRenderFor.Username)

	err := g.board.Render(w, playerToRenderFor, g)
	if err != nil {
		logger.Gamesf("Failed to render for %s: %v", playerToRenderFor.Username, err)
		return err
	}

	flusher.Flush()

	for {
		select {
		case <-playerToRenderFor.Events:
			err := g.board.Render(w, playerToRenderFor, g)
			if err != nil {
				logger.Gamesf("Failed to render for %s: %v", playerToRenderFor.Username, err)
				return err
			}

			flusher.Flush()
		case <-g.done:
			return nil
		}
	}
}
