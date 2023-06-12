package pollen

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func TestGame(t *testing.T) {
	g := NewGame(uuid.Must(uuid.NewV4()), "JAMES", "test game", 1)
	g.AddPlayer("RAE")
	require.NoError(t, g.Start())

	err := g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, Position{0.5, 0.5})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("RAE", g.GetHand("RAE")[0].ID, Position{0.5, -0.5})
	require.NoError(t, err)

	tk := g.GetNextToken()

	err = g.PlayToken("RAE", tk, Position{1, 0})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, Position{-0.5, 0.5})
	require.NoError(t, err)

	tk = g.GetNextToken()

	err = g.PlayToken("JAMES", tk, Position{0, 1})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("RAE", g.GetHand("RAE")[0].ID, Position{-0.5, -0.5})
	require.NoError(t, err)

	tk = g.GetNextToken()

	err = g.PlayToken("RAE", tk, Position{0, -1})
	require.NoError(t, err)

	tk = g.GetNextToken()

	err = g.PlayToken("RAE", tk, Position{-1, 0})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, Position{-1.5, -0.5})
	require.NoError(t, err)

	tk = g.GetNextToken()

	err = g.PlayToken("JAMES", tk, Position{-1, -1})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("RAE", g.GetHand("RAE")[0].ID, Position{-1.5, -1.5})
	require.NoError(t, err)

	tk = g.GetNextToken()

	err = g.PlayToken("RAE", tk, Position{-2, -1})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	// err = g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, Position{-1.5, -1.5})
	// if err != nil {
	// 	t.Fatal(err)
	// }
}

func TestMustPlayToken(t *testing.T) {
	g := NewGame(uuid.Must(uuid.NewV4()), "JAMES", "test 2 game", 1)
	g.AddPlayer("RAE")
	require.NoError(t, g.Start())

	err := g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, Position{0.5, 0.5})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("RAE", g.GetHand("RAE")[0].ID, Position{0.5, -0.5})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.Error(t, err)
}
