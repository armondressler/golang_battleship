package player

import (
	"github.com/google/uuid"
)

var Players []Player

type Player struct {
	name string
	id   uuid.UUID
}

func (player *Player) String() string {
	return player.name
}

func NewPlayer(name string) Player {
	id := uuid.New()
	p := Player{name, id}
	Players = append(Players, p)
	return p
}
