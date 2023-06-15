package pollen

import (
	"errors"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/logger"
)

var (
	ErrNoCard                    = errors.New("player has no card")
	ErrCardAlreadyPlayedThisTurn = errors.New("card already played this turn")
)

type Player struct {
	Color      Color
	Username   string
	Events     chan struct{}
	Hand       []GardenCard
	Deck       *GardenDeck
	HintsOn    bool
	cardPlayed bool
	Connected  bool
	l          sync.Mutex
}

func NewPlayer(username string, numPlayers int, color Color) *Player {
	p := &Player{
		Color:    color,
		Username: username,
		Events:   make(chan struct{}, 10),
		Deck:     NewGardenDeck(numPlayers, color),
	}
	p.Hand = p.Deck.cards[:5]
	p.Deck.cards = p.Deck.cards[5:]
	return p
}

func (p *Player) ToggleConnection() {
	p.l.Lock()
	p.Connected = !p.Connected
	logger.Usersf("Player %s is now connected: %v", p.Username, p.Connected)
	p.l.Unlock()
}

func (p *Player) IsConnected() bool {
	p.l.Lock()
	c := p.Connected
	p.l.Unlock()
	return c
}

func (p *Player) OutOfCards() bool {
	return len(p.Hand) == 0
}

func (p *Player) CardNotPlayed() {
	p.cardPlayed = false
}

func (p *Player) CardPlayed() bool {
	return p.cardPlayed
}

func (p *Player) PlayCard(card uuid.UUID) (*GardenCard, error) {
	if p.cardPlayed {
		return nil, ErrCardAlreadyPlayedThisTurn
	}
	if len(p.Hand) == 0 {
		return nil, ErrNoCard
	}
	for i, c := range p.Hand {
		if c.ID == card {
			logger.Player("Valid move ", c.Color, c.Type, c.Value)
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			drawnCard := p.Deck.Draw()
			if drawnCard != nil {
				p.Hand = append(p.Hand, *drawnCard)
			}
			p.cardPlayed = true
			return &c, nil
		}
	}
	return nil, errors.New("card not found")
}
func (p *Player) GetCard(card uuid.UUID) (*GardenCard, error) {
	if len(p.Hand) == 0 {
		return nil, ErrNoCard
	}
	for _, c := range p.Hand {
		if c.ID == card {
			return &c, nil
		}
	}
	return nil, errors.New("card not found")
}
