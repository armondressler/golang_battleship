package board

import (
	"fmt"
	"golang_battleship/ship"
	"testing"
)

func TestDeployment(t *testing.T) {
	someShips := []ship.Ship{*ship.NewShip("Destroyer", 3, 3, "n")}
	aBoard := Board{8, 8, someShips, 5, 2}
	fmt.Println(aBoard)
	fmt.Println(aBoard.ships)
}
