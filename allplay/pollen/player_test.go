package pollen

import (
	"testing"

	"github.com/hibooboo2/ggames/allplay/pollen/colors"
	"github.com/stretchr/testify/require"
)

func TestPlayerEmptyDeckAndHand(t *testing.T) {
	p := NewPlayer("JAMES", 2, colors.Green)

	card := p.Hand[0]
	for len(p.Hand) > 0 {
		card = p.Hand[0]
		_, err := p.PlayCard(card.ID)
		require.NoError(t, err)
		p.CardNotPlayed()
	}
	_, err := p.PlayCard(card.ID)
	require.Error(t, err)
}
