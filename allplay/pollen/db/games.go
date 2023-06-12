package db

import (
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/pollen"
)

var games = map[uuid.UUID]*pollen.Game{}

func NewGame(id uuid.UUID, username string) error {
	if _, ok := games[id]; ok {
		return errors.New("game already exists")
	}

	g := pollen.NewGame(id, username, "Like you've never danced before:", len(games)+1)
	g.AutoToken = true
	g.AddPlayer(username)
	games[id] = g
	return nil
}
func GetGame(id uuid.UUID) *pollen.Game {
	return games[id]
}

func GetInvitedGames(username string) []*pollen.Game {
	invitedGames := []*pollen.Game{}
	for _, game := range games {
		if game.HasInvitedUser(username) {
			invitedGames = append(invitedGames, game)
		}
	}
	return invitedGames
}

func GetActiveGames(username string) []*pollen.Game {
	activeGames := []*pollen.Game{}
	for _, game := range games {
		if game.HasPlayer(username) {
			activeGames = append(activeGames, game)
		}
	}
	return activeGames
}

func AddGameUser(gameID uuid.UUID, username string) error {
	g, ok := games[gameID]
	if !ok {
		return errors.New("game not found")
	}

	err := g.AddPlayer(username)
	if err != nil {
		return fmt.Errorf("failed to add player: %q %w", username, err)
	}

	return nil
}

func StartGame(gameID uuid.UUID) error {
	g, ok := games[gameID]
	if !ok {
		return errors.New("game not found")
	}
	err := g.Start()
	if err != nil {
		return fmt.Errorf("failed to start game: %q %w", gameID, err)
	}

	return nil
}

func InviteUserToGame(gameID uuid.UUID, username string) error {
	g, ok := games[gameID]
	if !ok {
		return errors.New("game not found")
	}

	if g.Started() {
		return errors.New("game already started")
	}
	g.InvitePlayer(username)
	return nil
}
