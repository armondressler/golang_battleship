package ship

import (
	"testing"
)

func TestCollides(t *testing.T) {
	shipOne := NewShip("Destroyer", 4, 5, "e")
	shipTwo := NewShip("Frigate", 5, 4, "n")
	shipThree := NewShip("Cruiser", 6, 3, "n")
	if !shipOne.Collides(*shipTwo) || !shipOne.Collides(*shipThree) || shipTwo.Collides(*shipThree) {
		t.Fail()
	}

}

func TestBowCoordinate(t *testing.T) {
	shipOne := NewShip("Destroyer", 4, 5, "e")
	if shipOne.BowCoordinate().x != 7 || shipOne.BowCoordinate().y != 5 {
		t.Fail()
	}

	shipTwo := NewShip("Frigate", 5, 4, "n")
	if shipTwo.BowCoordinate().x != 5 || shipTwo.BowCoordinate().y != 6 {
		t.Fail()
	}
}

func TestSternCoordinate(t *testing.T) {
	shipOne := NewShip("Destroyer", 4, 5, "e")
	if shipOne.SternCoordinate().x != 4 || shipOne.SternCoordinate().y != 5 {
		t.Fail()
	}
	shipTwo := NewShip("Frigate", 3, 3, "s")
	if shipTwo.SternCoordinate().x != 3 || shipTwo.SternCoordinate().y != 3 {
		t.Fail()
	}
}

func TestOrientation(t *testing.T) {
	shipOne := NewShip("Destroyer", 4, 5, "e")
	if shipOne.Orientation().x != 1 || shipOne.Orientation().y != 0 {
		t.Fail()
	}
	shipTwo := NewShip("Frigate", 5, 4, "n")
	if shipTwo.Orientation().x != 0 || shipTwo.Orientation().y != 1 {
		t.Fail()
	}
}
