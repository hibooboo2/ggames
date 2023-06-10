package pollen

import "github.com/gofrs/uuid"

type PlayerInput struct {
	Username string
	Input    chan PlayerCommand
}

type PlayerCommand struct {
	Command string
	Args    []uuid.UUID
}
