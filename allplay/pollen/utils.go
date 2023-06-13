package pollen

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

func GetGameID(r *http.Request) uuid.UUID {
	id := chi.URLParam(r, "game_id")
	gameID := uuid.Must(uuid.FromString(id))
	return gameID
}

func GetCardID(r *http.Request) uuid.UUID {
	id := chi.URLParam(r, "card_id")
	cardID := uuid.Must(uuid.FromString(id))
	return cardID
}

func GetTokenID(r *http.Request) uuid.UUID {
	id := chi.URLParam(r, "token_id")
	tokenID := uuid.Must(uuid.FromString(id))
	return tokenID
}

func GetPosition(r *http.Request) Position {
	positionString := r.FormValue("position")
	if positionString == "" {
		panic("position is empty")
	}
	position, err := ParsePosition(positionString)
	if err != nil {
		panic(err)
	}
	return *position
}
