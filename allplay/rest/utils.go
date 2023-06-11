package rest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/pollen"
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

func GetPosition(r *http.Request) pollen.Position {
	positionString := r.FormValue("position")
	if positionString == "" {
		panic("position is empty")
	}
	position, err := pollen.ParsePosition(positionString)
	if err != nil {
		panic(err)
	}
	return *position
}
