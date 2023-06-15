package pollen

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/pollen/cards"
	"github.com/hibooboo2/ggames/allplay/pollen/position"
	"github.com/hibooboo2/ggames/allplay/pollen/token"
	"github.com/stretchr/testify/require"
)

func TestScoring(t *testing.T) {
	//XXX make token bag get token able to not be shuffled and also make it be consistent if got twice with no consume

	// Create a board and play tokens and cards on it then check scores
	// Should be able to ignore turns etc.. Since its just for scores and board
	bag := token.NewTokenBag(false)
	tk := bag.TakeNextToken()

	g := NewGame(uuid.Must(uuid.NewV4()), "JAMES", "Some Game")
	b := NewBoard(tk, g)

	require.False(t, b.CanPlayToken(position.Position{-1, 0}) == nil)

	dk1 := cards.NewGardenDeck(2, 0)
	dk2 := cards.NewGardenDeck(2, 1)

	require.NoError(t, b.PlayCard(position.Position{-0.5, 0.5}, dk1.Draw()))
	require.NoError(t, b.PlayCard(position.Position{0.5, -0.5}, dk2.Draw()))

	require.NoError(t, b.PlayCard(position.Position{0.5, 0.5}, dk1.Draw()))

	require.NoError(t, b.PlayToken(position.Position{0, 1}, bag.TakeNextToken()))
	require.NoError(t, b.PlayToken(position.Position{1, 0}, bag.TakeNextToken()))

	require.NoError(t, b.PlayCard(position.Position{-0.5, -0.5}, dk1.Draw()))

	require.NoError(t, b.PlayToken(position.Position{0, -1}, bag.TakeNextToken()))
	require.NoError(t, b.PlayToken(position.Position{-1, 0}, bag.TakeNextToken()))

	b.UpdateScores()

	require.Equal(t, 10, b.Scores.BeeMeeplesLeft)
	require.Equal(t, 10, b.Scores.JunebugMeeplesLeft)
	require.Equal(t, 10, b.Scores.ButterflyMeeplesLeft)

}
