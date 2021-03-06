package board

import (
	"bytes"
	"fmt"
	"golang_battleship/ship"
	"golang_battleship/weapon"
	"io"
)

type Board struct {
	BoardParameters
	ships   []ship.Ship
	impacts []impact
}

type impact struct {
	x, y   int
	weapon weapon.Exploder
}

type coordinate struct {
	x, y int
}

type BoardParameters struct {
	SizeX    int `json:"size_x"`
	SizeY    int `json:"size_y"`
	MaxShips int `json:"max_ships"`
}

func NewBoard(bp BoardParameters) Board {
	return Board{bp, []ship.Ship{}, []impact{}}
}

func (board Board) Size() int {
	return board.BoardParameters.SizeX * board.BoardParameters.SizeY
}

func (board Board) Dimensions() (int, int) {
	return board.BoardParameters.SizeX, board.BoardParameters.SizeY
}

func (board Board) String() string {
	var buf bytes.Buffer
	board.draw(&buf)
	return buf.String()
}

func (board Board) draw(writer io.Writer) {
	shipArray := board.unpackShips()
	impactArray := board.unpackImpacts()
	for y := board.BoardParameters.SizeX; y >= 0; y-- {
		for x := 0; x < int(board.BoardParameters.SizeY); x++ {
			if symbol, ok := impactArray[coordinate{x, y}]; ok {
				writer.Write([]byte(string(symbol)))
			} else if symbol, ok := shipArray[coordinate{x, y}]; ok {
				writer.Write([]byte(string(symbol)))
			} else {
				writer.Write([]byte(string('#')))
			}
			writer.Write([]byte(string(' ')))
		}
		io.WriteString(writer, "\n")
	}
}

func (board Board) unpackShips() map[coordinate]rune {
	shipArray := make(map[coordinate]rune)
	for _, ship := range board.ships {
		for _, c := range ship.Coordinates() {
			shipArray[coordinate{c.X(), c.Y()}] = ship.Symbol()
		}
	}
	return shipArray
}

func (board Board) unpackImpacts() map[coordinate]rune {
	impactArray := make(map[coordinate]rune)
	for _, impact := range board.impacts {
		impactArray[coordinate{impact.x, impact.y}] = impact.weapon.Symbol()
	}
	return impactArray
}

func (board Board) checkCollision(ship ship.Ship) *ship.Ship {
	for _, otherShip := range board.ships {
		if otherShip.Collides(ship) {
			return &otherShip
		}
	}
	return nil
}

func (board *Board) DeployShip(ship ship.Ship) error {
	if collidingShip := board.checkCollision(ship); collidingShip != nil {
		return fmt.Errorf("collision with ship %s detected", collidingShip)
	}
	board.ships = append(board.ships, ship)
	return nil
}
