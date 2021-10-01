package player

type Player struct {
	name string
	id   int64
}

func (player *Player) String() string {
	return player.name
}
