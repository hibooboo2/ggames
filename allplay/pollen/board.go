package pollen

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"math"

	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/logger"
	"github.com/hibooboo2/ggames/allplay/pollen/cards"
	"github.com/hibooboo2/ggames/allplay/pollen/position"
	"github.com/hibooboo2/ggames/allplay/pollen/token"
)

type Board struct {
	cards   map[position.Position]*cards.GardenCard
	tokens  map[position.Position]*token.PollinatorToken
	Scores  *GameScore
	Players int
	g       *Game
}

func NewBoard(tk *token.PollinatorToken, g *Game) *Board {
	players := len(g.Players)
	b := &Board{
		cards:   make(map[position.Position]*cards.GardenCard),
		tokens:  make(map[position.Position]*token.PollinatorToken),
		Players: players,
		Scores:  NewGameScore(players),
		g:       g,
	}
	err := b.PlaceStarterToken(tk)
	if err != nil {
		panic(err)
	}
	return b
}

func (b *Board) PlaceStarterToken(token *token.PollinatorToken) error {
	err := b.PlayToken(position.Position{0, 0}, token)
	if err != nil {
		return fmt.Errorf("failed to place starter token: %w", err)
	}
	return nil
}

func (b *Board) MustPlayToken() error {
	if b.GameOver() {
		return nil
	}

	if len(b.cards) == 0 {
		return nil
	}
	positions := b.GetTokensMustPlay()
	if len(positions) == 0 {
		return nil
	}
	return fmt.Errorf("must play tokens at the following positions: %v", positions)
}

func (b *Board) GetTokensMustPlay() []position.Position {
	if b.GameOver() {
		return nil
	}

	if len(b.cards) == 0 {
		return nil
	}
	tokensMustPlay := map[position.Position]struct{}{}
	for p, _ := range b.cards {
		nw := position.Position{p.X + 0.5, p.Y + 0.5}
		sw := position.Position{p.X + 0.5, p.Y - 0.5}
		ne := position.Position{p.X - 0.5, p.Y + 0.5}
		se := position.Position{p.X - 0.5, p.Y - 0.5}

		logger.AtLevel(logger.LPosition, "Checking nw position for position", p)
		if b.CanPlayToken(nw) == nil {
			tokensMustPlay[nw] = struct{}{}
		}
		logger.AtLevel(logger.LPosition, "Checking sw position for position", p)
		if b.CanPlayToken(sw) == nil {
			tokensMustPlay[sw] = struct{}{}
		}
		logger.AtLevel(logger.LPosition, "Checking ne position for position", p)
		if b.CanPlayToken(ne) == nil {
			tokensMustPlay[ne] = struct{}{}
		}
		logger.AtLevel(logger.LPosition, "Checking se position for position", p)
		if b.CanPlayToken(se) == nil {
			tokensMustPlay[se] = struct{}{}
		}
	}
	var positions []position.Position
	for position := range tokensMustPlay {
		positions = append(positions, position)
	}
	logger.AtLevel(logger.LPosition, "Found %v tokens to play", len(positions))
	return positions
}

func (b *Board) PlayToken(p position.Position, token *token.PollinatorToken) error {
	if token.IsUsed() {
		return fmt.Errorf("token is already used: %v", token.ID)
	}
	_, ok := b.tokens[p]
	if ok {
		return fmt.Errorf("token already played at position %v", p)
	}

	logger.AtLevelf(logger.LBoard|logger.LToken|logger.LPosition, "Trying to place token at position %v id: %s", p, token.ID)
	if token.Position != nil && p.X != 0 && p.Y != 0 {
		logger.AtLevelf(logger.LBoard|logger.LToken|logger.LPosition, "Token already has a position %v cannot play at: %v id: %s", *token.Position, p, token.ID)
		return fmt.Errorf("token already played at position %v", *token.Position)
	}

	err := b.CanPlayToken(p)
	if err != nil {
		return fmt.Errorf("cannot play token: %w", err)
	}
	token.Position = &p
	b.tokens[p] = token
	token.Play()
	return nil
}

func (b *Board) CanPlayToken(p position.Position) error {
	if _, present := b.tokens[p]; present {
		return fmt.Errorf("token already exists at position %v", p)
	}

	_, fracX := math.Modf(math.Abs(p.X))
	_, fracY := math.Modf(math.Abs(p.Y))
	logger.AtLevel(logger.LPosition, "Checking if can play: ", p, fracX, fracY)
	if fracX != 0 || fracY != 0 {
		return fmt.Errorf("invalid position for a token: %v X and Y positions must be whole numbers", p)
	}

	if p.X == 0 && p.Y == 0 {
		return nil
	}

	//Add verification that position is a whole number
	ne := position.Position{p.X + 0.5, p.Y + 0.5}
	se := position.Position{p.X + 0.5, p.Y - 0.5}
	nw := position.Position{p.X - 0.5, p.Y + 0.5}
	sw := position.Position{p.X - 0.5, p.Y - 0.5}
	_, swPresent := b.cards[sw]
	_, nwPresent := b.cards[nw]
	_, nePresent := b.cards[ne]
	_, sePresent := b.cards[se]
	switch {
	case swPresent && nwPresent:
	case sePresent && nePresent:
	case swPresent && sePresent:
	case nwPresent && nePresent:
	default:
		return fmt.Errorf("there is not two adjacent cards at position %v", p)
	}
	return nil
}

func (b *Board) TokenIsSurrounded(tokenID uuid.UUID) bool {
	var p *position.Position
	for p2, token := range b.tokens {
		if token.ID == tokenID {
			if token.IsSurrounded() {
				return true
			}
			p = &p2
			break
		}
	}

	if p == nil {
		return false
	}

	logger.AtLevel(logger.LBoard|logger.LPosition, logger.LToken, "Checking if token is surrounded: ", p)

	//Add verification that position is a whole number
	ne := position.Position{p.X + 0.5, p.Y + 0.5}
	se := position.Position{p.X + 0.5, p.Y - 0.5}
	nw := position.Position{p.X - 0.5, p.Y + 0.5}
	sw := position.Position{p.X - 0.5, p.Y - 0.5}
	_, swPresent := b.cards[sw]
	_, nwPresent := b.cards[nw]
	_, nePresent := b.cards[ne]
	_, sePresent := b.cards[se]
	isSurrounded := swPresent && nwPresent && nePresent && sePresent

	logger.AtLevelf(logger.LBoard|logger.LPosition|logger.LScore, "Token is directions sw: %v se %v nw %v ne %v", swPresent, sePresent, nwPresent, nePresent)
	logger.AtLevelf(logger.LBoard|logger.LPosition|logger.LScore, "Token is surrounded: %v %v", *p, isSurrounded)
	return isSurrounded
}

type MeepleType int

const (
	BeeMeeple MeepleType = iota
	JunebugMeeple
	ButterflyMeeple
)

func (b *Board) PlayCard(p position.Position, card *cards.GardenCard) error {
	err := b.CanPlayCard(p)
	if err != nil {
		return fmt.Errorf("cannot play card: %w", err)
	}
	b.cards[p] = card
	b.UpdateScores()
	return nil
}

func (b *Board) CanPlayCard(p position.Position) error {
	if err := b.MustPlayToken(); err != nil {
		return fmt.Errorf("you must play a token first: %w", err)
	}
	_, fracX := math.Modf(math.Abs(p.X))
	_, fracY := math.Modf(math.Abs(p.Y))
	if fracX != 0.5 || fracY != 0.5 {
		return fmt.Errorf("invalid position for a card: %v X and Y positions must be numbers that end in .5", p)
	}

	_, present := b.cards[p]
	if present {
		return fmt.Errorf("card already exists at position %v", p)
	}

	//Add verification that position is a Half number
	ne := position.Position{p.X + 0.5, p.Y + 0.5}
	se := position.Position{p.X + 0.5, p.Y - 0.5}
	nw := position.Position{p.X - 0.5, p.Y + 0.5}
	sw := position.Position{p.X - 0.5, p.Y - 0.5}
	_, swPresent := b.tokens[sw]
	_, nwPresent := b.tokens[nw]
	_, nePresent := b.tokens[ne]
	_, sePresent := b.tokens[se]
	switch {
	case swPresent, nwPresent, nePresent, sePresent:
	default:
		return fmt.Errorf("position %v does not have an adjacent token", p)
	}
	return nil
}

func (b *Board) CardLocationsPlayable() map[position.Position]struct{} {
	if b.GameOver() {
		return nil
	}

	positions := map[position.Position]struct{}{}
	for p := range b.tokens {
		nw := position.Position{p.X + 0.5, p.Y + 0.5}
		sw := position.Position{p.X + 0.5, p.Y - 0.5}
		ne := position.Position{p.X - 0.5, p.Y + 0.5}
		se := position.Position{p.X - 0.5, p.Y - 0.5}

		if b.CanPlayCard(nw) == nil {
			positions[nw] = struct{}{}
		}
		if b.CanPlayCard(sw) == nil {
			positions[sw] = struct{}{}
		}
		if b.CanPlayCard(ne) == nil {
			positions[ne] = struct{}{}
		}
		if b.CanPlayCard(se) == nil {
			positions[se] = struct{}{}
		}
	}
	return positions
}

func (b *Board) GameOver() bool {
	return b.g.AllPlayersOutOfCards() || b.Scores.OutOfMeeples() || b.g.TokenBag.OutOfTokens()
}

var boardFuncs = template.FuncMap{
	"offset": func(v float64, o int) int {
		return int(v*2*25) + o
	},
}

func (b *Board) Render(w io.Writer, p *Player, g *Game) error {
	boardTmpl, err := LoadTemplate("./pollen/static/views/board.html.tmpl", boardFuncs)
	if err != nil {
		return fmt.Errorf("board template not found: %w", err)
	}

	buff := bytes.NewBuffer(nil)
	tokensMustPlay := b.GetTokensMustPlay()
	nextTokens := len(tokensMustPlay)
	if nextTokens == 0 {
		nextTokens = 1
	}
	err = boardTmpl.ExecuteTemplate(buff, "board", struct {
		Cards                  map[position.Position]*cards.GardenCard
		Tokens                 map[position.Position]*token.PollinatorToken
		PlayableCards          map[position.Position]struct{}
		PlayableTokenPositions []position.Position
		TokensCanPlay          []*token.PollinatorToken
		Debug                  bool
		Player                 *Player
		GameID                 string
		Hand                   []cards.GardenCard
		IsPlayerTurn           bool
		HintsOn                bool
		Scores                 *GameScore
		GameOver               bool
		Players                map[string]*Player
	}{
		Cards:                  b.cards,
		Tokens:                 b.tokens,
		PlayableCards:          b.CardLocationsPlayable(),
		PlayableTokenPositions: tokensMustPlay,
		TokensCanPlay:          g.TokenBag.GetTokens(nextTokens),
		Debug:                  false,
		Player:                 p,
		GameID:                 g.id.String(),
		Hand:                   p.Hand,
		IsPlayerTurn:           g.activePlayer().Username == p.Username,
		HintsOn:                p.HintsOn,
		Scores:                 b.Scores,
		GameOver:               b.GameOver(),
		Players:                g.Players,
	})
	if err != nil {
		return err
	}

	logger.Boardln("Rendering board:", p.Color)
	fmt.Fprintf(w, "data: %s\n\n", base64.StdEncoding.EncodeToString(buff.Bytes()))
	return nil
}
