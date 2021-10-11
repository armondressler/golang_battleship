package game

import (
	"golang_battleship/board"
	"golang_battleship/player"
	"golang_battleship/ship"

	"github.com/google/uuid"
)

type game struct {
	participants []participant
	id           uuid.UUID
}

type participant struct {
	player player.Player
	ship   []ship.Ship
	board  board.Board
	id     uuid.UUID
}
