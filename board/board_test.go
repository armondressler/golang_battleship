package board

import (
	"fmt"
	"golang_battleship/ship"
	"testing"
)

func TestDeployment(t *testing.T) {
	someShips := []ship.Ship{*ship.NewShip("Destroyer", 3, 3, "n")}
	aBoard := Board{BoardParameters{8, 8, 4}, someShips, []impact{}}
	fmt.Println(aBoard)
	fmt.Println(aBoard.ships)
}

func TestDraw(t *testing.T) {
	someShips := []ship.Ship{
		*ship.NewShip("Destroyer", 0, 0, "n"),
		*ship.NewShip("Carrier", 1, 0, "e"),
		*ship.NewShip("Frigate", 5, 6, "s"),
	}
	aBoard := Board{BoardParameters{8, 8, 4}, someShips, []impact{}}
	drawThis := `# # # # # # # # 
# # # # # # # # 
# # # # # F # # 
# # # # # F # # 
# # # # # F # # 
D # # # # # # # 
D # # # # # # # 
D # # # # # # # 
D T T T T T T T 
`
	if aBoard.String() != drawThis {
		t.Fail()
	}
}
