package player

import (
	"github.com/google/uuid"
)

type Player struct {
	name string
	id   uuid.UUID
}

func (player *Player) String() string {
	return player.name
}

func NewPlayer(name string) Player {
	id := uuid.New()
	return Player{name, id}
}
