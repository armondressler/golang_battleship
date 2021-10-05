package board

import (
	"bytes"
	"fmt"
	"golang_battleship/ship"
	"testing"
)

func TestDeployment(t *testing.T) {
	someShips := []ship.Ship{
		*ship.NewShip("Destroyer", 3, 3, "n"),
		*ship.NewShip("Frigate", 5, 6, "s"),
	}
	aBoard := Board{8, 8, someShips, []impact{}, 4}
	fmt.Println(len(aBoard.ships))
	if len(aBoard.ships) != 2 {
		t.Fail()
	}
	oneMoreShipThatCollides := ship.NewShip("Frigate", 5, 6, "s")
	if err := aBoard.DeployShip(*oneMoreShipThatCollides); err == nil {
		t.Fail()
	}
}

func TestDraw(t *testing.T) {
	someShips := []ship.Ship{
		*ship.NewShip("Destroyer", 0, 0, "n"),
		*ship.NewShip("Carrier", 1, 0, "e"),
		*ship.NewShip("Frigate", 5, 6, "s"),
	}
	aBoard := Board{8, 8, someShips, []impact{}, 4}
	var buf bytes.Buffer
	aBoard.draw(&buf)
	fmt.Println(buf.String())
	fmt.Println(len(buf.String()))
}
