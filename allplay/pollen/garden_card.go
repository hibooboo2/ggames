package pollen

import (
	"fmt"
	"math/rand"

	"github.com/gofrs/uuid"
)

type GardenCard struct {
	ID    uuid.UUID
	Type  CardType
	Value int
	Color Color
}

func (c GardenCard) Name() string {
	return fmt.Sprintf("%s_%s_%d", c.Type, c.Color, c.Value)
}

type GardenDeck struct {
	cards []GardenCard
}

func NewGardenDeck(numplayers int, color Color) *GardenDeck {
	deck := append(gardenCardSuits(color), gardenCardWilds(numplayers, color)...)
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return &GardenDeck{deck}
}

func (d *GardenDeck) Draw() *GardenCard {
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

func gardenCardSuits(color Color) []GardenCard {
	return []GardenCard{
		{uuid.Must(uuid.NewV4()), Bee, 2, color},
		{uuid.Must(uuid.NewV4()), Bee, 3, color},
		{uuid.Must(uuid.NewV4()), Bee, 4, color},
		{uuid.Must(uuid.NewV4()), Butterfly, 2, color},
		{uuid.Must(uuid.NewV4()), Butterfly, 3, color},
		{uuid.Must(uuid.NewV4()), Butterfly, 4, color},
		{uuid.Must(uuid.NewV4()), Junebug, 2, color},
		{uuid.Must(uuid.NewV4()), Junebug, 3, color},
		{uuid.Must(uuid.NewV4()), Junebug, 4, color},
	}
}

func gardenCardWilds(numplayers int, color Color) []GardenCard {
	cards := []GardenCard{
		{uuid.Must(uuid.NewV4()), Wild, 1, color},
		{uuid.Must(uuid.NewV4()), Wild, 1, color},
		{uuid.Must(uuid.NewV4()), Wild, 2, color},
		{uuid.Must(uuid.NewV4()), Wild, 3, color},
	}
	switch numplayers {
	case 2, 3:
		return append(cards, GardenCard{uuid.Must(uuid.NewV4()), Wild, 1, color}, GardenCard{uuid.Must(uuid.NewV4()), Wild, 2, color})
	case 4:
		return cards
	default:
		panic("invalid number of players")
	}
}
