package pollen

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/logger"
	"github.com/hibooboo2/ggames/allplay/pollen/position"
	"github.com/hibooboo2/glog"
	"github.com/stretchr/testify/require"
)

func init() {
	logger.SetLevel(logger.LAuth | logger.LGames | logger.LPlayer | logger.LUsers | glog.DefaultLevel | logger.LInit)
}

func TestGame(t *testing.T) {
	g := NewGame(uuid.Must(uuid.NewV4()), "JAMES", "test game")
	g.AddPlayer("RAE")
	require.NoError(t, g.Start())

	err := g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, position.Position{0.5, 0.5})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("RAE", g.GetHand("RAE")[0].ID, position.Position{0.5, -0.5})
	require.NoError(t, err)

	tk := g.GetNextTokenID()

	err = g.PlayToken("RAE", *tk, position.Position{1, 0})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, position.Position{-0.5, 0.5})
	require.NoError(t, err)

	tk = g.GetNextTokenID()

	err = g.PlayToken("JAMES", *tk, position.Position{0, 1})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("RAE", g.GetHand("RAE")[0].ID, position.Position{-0.5, -0.5})
	require.NoError(t, err)

	tk = g.GetNextTokenID()

	err = g.PlayToken("RAE", *tk, position.Position{0, -1})
	require.NoError(t, err)

	tk = g.GetNextTokenID()

	err = g.PlayToken("RAE", *tk, position.Position{-1, 0})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, position.Position{-1.5, -0.5})
	require.NoError(t, err)

	tk = g.GetNextTokenID()

	err = g.PlayToken("JAMES", *tk, position.Position{-1, -1})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("RAE", g.GetHand("RAE")[0].ID, position.Position{-1.5, -1.5})
	require.NoError(t, err)

	tk = g.GetNextTokenID()

	err = g.PlayToken("RAE", *tk, position.Position{-2, -1})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	// err = g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, position.Position{-1.5, -1.5})
	// if err != nil {
	// 	t.Fatal(err)
	// }
}

func TestMustPlayToken(t *testing.T) {
	g := NewGame(uuid.Must(uuid.NewV4()), "JAMES", "test 2 game")
	g.AddPlayer("RAE")
	require.NoError(t, g.Start())

	err := g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, position.Position{0.5, 0.5})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("RAE", g.GetHand("RAE")[0].ID, position.Position{-0.5, -0.5})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.NoError(t, err)

	err = g.PlayCard("JAMES", g.GetHand("JAMES")[0].ID, position.Position{0.5, -0.5})
	require.NoError(t, err)

	err = g.NextPlayer()
	require.Error(t, err)
}
