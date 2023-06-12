package pollen

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/gofrs/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Game struct {
	id                 uuid.UUID
	players            []*Player
	playerUsernames    map[string]struct{}
	activePlayerCursor int
	tokenBag           *TokenBag
	board              *Board
	events             chan struct{}
	done               chan struct{}
	started            bool
}

func NewGame(id uuid.UUID) *Game {
	g := &Game{
		id:              id,
		tokenBag:        NewTokenBag(),
		events:          make(chan struct{}),
		done:            make(chan struct{}),
		playerUsernames: map[string]struct{}{},
	}
	return g
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
		log.Println(p.Username, p.Color)
	}
	g.board = NewBoard(g.tokenBag.GetToken())
	go func() {
		for range g.events {
			for _, p := range g.players {
				select {
				case p.Events <- struct{}{}:
				case <-time.After(time.Millisecond * 10):
				}
			}
		}
	}()
	g.started = true
	return nil
}

func (g *Game) AddPlayer(username string) error {
	if len(g.playerUsernames) > 3 {
		return errors.New("game already has 4 players")
	}

	log.Printf("Adding user to game: %s %s", g.id, username)
	g.playerUsernames[username] = struct{}{}
	return nil
}

func (g *Game) HasPlayer(username string) bool {
	for uname := range g.playerUsernames {
		log.Printf("Checking if user is in game: %s %q", g.id, uname)
		if uname == username {
			return true
		}
	}
	return false
}

func (g *Game) End() {
	close(g.done)
}

func (g *Game) GetID() uuid.UUID {
	return g.id
}

func (g *Game) GetHand(username string) []GardenCard {
	log.Println("Getting hand for", username)
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

	g.events <- struct{}{}
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
	if playerToRenderFor == nil {
		return fmt.Errorf("player %q not found", username)
	}

	log.Println("Starting render for", playerToRenderFor.Username)
	t := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-t.C:
		case <-playerToRenderFor.Events:
		case <-g.done:
			return nil
		}

		err := g.board.Render(w, playerToRenderFor, g)
		if err != nil {
			log.Printf("Failed to render for %s: %v", playerToRenderFor.Username, err)
			return err
		}
		switch w := w.(type) {
		case http.Flusher:
			w.Flush()
		}
	}
}
