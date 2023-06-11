package pollen

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlayerEmptyDeckAndHand(t *testing.T) {
	p := NewPlayer("JAMES", 2, Green)

	card := p.Hand[0]
	for len(p.Hand) > 0 {
		card = p.Hand[0]
		_, err := p.PlayCard(card.ID)
		require.NoError(t, err)
	}
	_, err := p.PlayCard(card.ID)
	require.Error(t, err)
}
