package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/pollen"
)

var games = map[uuid.UUID]*pollen.Game{}
var pollenNewGameUsers = map[uuid.UUID][]string{}

func NewGame(id uuid.UUID, username string) error {
	_, ok := games[id]
	if ok {
		return errors.New("game already exists")
	}

	pollenNewGameUsers[id] = []string{}
	return nil
}
func GetGame(id uuid.UUID) *pollen.Game {
	return games[id]
}

func AddGameUser(gameID uuid.UUID, username string) error {
	users, ok := pollenNewGameUsers[gameID]
	if !ok {
		return errors.New("game not found")
	}

	if len(users) > 3 {
		return errors.New("max users in game already reached")
	}

	for _, user := range users {
		if user == username {
			log.Println("user already in game")
			return nil
		}
	}

	log.Printf("Adding user to game: %s %s", gameID, username)
	pollenNewGameUsers[gameID] = append(pollenNewGameUsers[gameID], username)
	return nil
}

func StartGame(gameID uuid.UUID) error {
	users, ok := pollenNewGameUsers[gameID]
	if !ok {
		return errors.New("game not found")
	}
	if len(users) == 0 {
		return errors.New("no users in game")
	}
	if len(users) < 2 || len(users) > 4 {
		return fmt.Errorf("invalid number of users in game current users: %v", users)
	}

	games[gameID] = pollen.NewGame(gameID, users)
	return nil
}
