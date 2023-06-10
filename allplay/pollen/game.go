package pollen

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/gofrs/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewGame(users []PlayerInput) *Game {
	g := &Game{
		id:       uuid.Must(uuid.NewV4()),
		tokenBag: NewTokenBag(),
	}
	for i, user := range users {
		g.players = append(g.players, NewPlayer(user.Username, len(users), Color(i)))
	}
	for _, p := range g.players {
		log.Println(p.Username, p.Color)
	}
	g.board = NewBoard(g.tokenBag.GetToken())
	return g
}

type Game struct {
	id                 uuid.UUID
	players            []*Player
	activePlayerCursor int
	tokenBag           *TokenBag
	board              *Board
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
	return nil
}

func (g *Game) PlayCard(username string, card uuid.UUID, position Position) error {
	p := g.activePlayer()
	if p.Username != username {
		return fmt.Errorf("it is not %q turn", username)
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

	return nil
}

func (g *Game) Render(w io.Writer) error {
	return g.board.Render(w)
}
