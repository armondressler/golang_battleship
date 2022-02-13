package game

import (
	"golang_battleship/player"
	"testing"
)

func TestNewGame(t *testing.T) {
	x := 12
	y := x
	maxships := 6
	p1, _ := player.NewPlayer("Rudolf", "")
	p2, _ := player.NewPlayer("Dagobert", "")
	NewGame(x, y, maxships, "New Game", 2, p1.Name, p2.Name)
}
