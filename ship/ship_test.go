package ship

import (
	"fmt"
	"testing"
)

func TestPrint(t *testing.T) {
	aShip := NewShip("Carrier", 3, 5, "n")
	fmt.Println(aShip)
	fmt.Println(aShip.Coordinates())
	// if v != 1.5 {
	// 	t.Error("Expected 1.5, got ", v)
	// }
}
