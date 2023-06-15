package cards

import (
	"fmt"
	"math/rand"

	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/logger"
	"github.com/hibooboo2/ggames/allplay/pollen/colors"
)

type GardenCard struct {
	ID    uuid.UUID
	Type  CardType
	Value int
	C     colors.Color
}

func (c GardenCard) Name() string {
	return fmt.Sprintf("%s_%s_%d", c.Type, c.C, c.Value)
}

type GardenDeck struct {
	cards []GardenCard
}

func NewGardenDeck(numplayers int, c colors.Color) *GardenDeck {
	logger.Cardsf("Making new deck with %d players color %s", numplayers, c)
	deck := append(gardenCardSuits(c), gardenCardWilds(numplayers, c)...)
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return &GardenDeck{deck}
}

func (d *GardenDeck) DrawHand() []GardenCard {
	h := d.cards[:5]
	d.cards = d.cards[5:]
	return h
}

func (d *GardenDeck) Draw() *GardenCard {
	logger.Cardsln("Drawing card")
	if len(d.cards) == 0 {
		return nil
	}
	card := d.cards[0]

	if len(d.cards) == 1 {
		d.cards = nil
		return &card
	}

	d.cards = d.cards[1:]
	return &card
}

func gardenCardSuits(c colors.Color) []GardenCard {
	return []GardenCard{
		{uuid.Must(uuid.NewV4()), Bee, 2, c},
		{uuid.Must(uuid.NewV4()), Bee, 3, c},
		{uuid.Must(uuid.NewV4()), Bee, 4, c},
		{uuid.Must(uuid.NewV4()), Butterfly, 2, c},
		{uuid.Must(uuid.NewV4()), Butterfly, 3, c},
		{uuid.Must(uuid.NewV4()), Butterfly, 4, c},
		{uuid.Must(uuid.NewV4()), Junebug, 2, c},
		{uuid.Must(uuid.NewV4()), Junebug, 3, c},
		{uuid.Must(uuid.NewV4()), Junebug, 4, c},
	}
}

func gardenCardWilds(numplayers int, c colors.Color) []GardenCard {
	cards := []GardenCard{
		{uuid.Must(uuid.NewV4()), Wild, 1, c},
		{uuid.Must(uuid.NewV4()), Wild, 1, c},
		{uuid.Must(uuid.NewV4()), Wild, 2, c},
		{uuid.Must(uuid.NewV4()), Wild, 3, c},
	}
	switch numplayers {
	case 2, 3:
		return append(cards, GardenCard{uuid.Must(uuid.NewV4()), Wild, 1, c}, GardenCard{uuid.Must(uuid.NewV4()), Wild, 2, c})
	case 4:
		return cards
	default:
		panic("invalid number of players")
	}
}
