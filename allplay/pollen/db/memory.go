package db

import (
	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/pollen"
)

var (
	games    = map[uuid.UUID]*pollen.Game{}
	users    = map[string][32]byte{}
	sessions = map[string]*UserSession{}
)
