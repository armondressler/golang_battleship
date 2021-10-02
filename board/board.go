package board

import (
	"fmt"
	"golang_battleship/ship"
	"io"
)

type Board struct {
	x, y              int16
	ships             []ship.Ship
	maxShipsPerPlayer int8
	maxPlayers        int8
}

func (board Board) size() int64 {
	return int64(board.x * board.y)
}

func (board Board) String() string {
	return "fubar"
}

func (board Board) draw(writer io.Writer) {

}

func (board Board) checkCollision(ship ship.Ship) *ship.Ship {
	for _, otherShip := range board.ships {
		if otherShip.Collides(ship) {
			return &otherShip
		}
	}
	return nil
}

func (board *Board) deployShip(ship ship.Ship) error {
	if collidingShip := board.checkCollision(ship); collidingShip != nil {
		return fmt.Errorf("collision with ship %s detected", collidingShip)
	}
	board.ships = append(board.ships, ship)
	return nil
}
