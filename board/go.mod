module board

go 1.16

replace golang_battleship/ship => ../ship

require (
	golang_battleship/player v0.0.0-00010101000000-000000000000
	golang_battleship/ship v0.0.0-00010101000000-000000000000
)

replace golang_battleship/player => ../player
