package ship

import (
	"fmt"
)

type orientation struct {
	x, y int8
}

type coordinate struct {
	x, y int
}

type structureUnit struct {
	c       coordinate
	healthy bool
}

type class struct {
	name   string
	length int
}

var lengthMap = map[string]int{
	"Submarine": 2,
	"Frigate":   3,
	"Destroyer": 4,
	"Cruiser":   5,
	"Carrier":   7,
}

var orientationMap = map[string]orientation{
	"n": {0, 1},
	"e": {1, 0},
	"s": {0, -1},
	"w": {-1, 0},
}

type Ship struct {
	class           class
	sternCoordinate coordinate
	orientation     orientation
	structure       []structureUnit
}

func NewShip(className string, x int, y int, orientation string) *Ship {
	s := &Ship{class{className, lengthMap[className]}, coordinate{x, y}, orientationMap[orientation], []structureUnit{}}
	for i := 0; i < s.Length(); i++ {
		x := s.sternCoordinate.x + i*int(s.orientation.x)
		y := s.sternCoordinate.y + i*int(s.orientation.y)
		s.structure = append(s.structure, structureUnit{c: coordinate{x, y}, healthy: true})
	}
	return s
}

func (coordinate coordinate) String() string {
	return fmt.Sprintf("x:%d/y:%d", coordinate.x, coordinate.y)
}

func (o orientation) String() string {
	var verboseOrientation string
	if o.x == 0 {
		if o.y == 1 {
			verboseOrientation = "North"
		} else {
			verboseOrientation = "South"
		}
	} else if o.x == 1 {
		verboseOrientation = "East"
	} else {
		verboseOrientation = "West"
	}
	return verboseOrientation
}

func (ship Ship) String() string {
	return fmt.Sprintf("%v (Stern: %v, Heading: %v, Length: %d, Hits: %d)", ship.class.name, ship.SternCoordinate().String(), ship.orientation.String(), ship.Length(), ship.Hits())
}

func (ship Ship) Length() int {
	return ship.class.length
}

func (ship Ship) Coordinates() []coordinate {
	var retval = []coordinate{}
	for i := 0; i < ship.Length(); i++ {
		x := ship.sternCoordinate.x + i*int(ship.orientation.x)
		y := ship.sternCoordinate.y + i*int(ship.orientation.y)
		retval = append(retval, coordinate{x, y})
	}
	return retval
}

func (ship Ship) BowCoordinate() coordinate {
	return coordinate{ship.sternCoordinate.x + ship.Length()*int(ship.orientation.x),
		ship.sternCoordinate.y + ship.Length()*int(ship.orientation.y)}
}

func (ship Ship) SternCoordinate() coordinate {
	return coordinate{ship.sternCoordinate.x, ship.sternCoordinate.y}
}

func (ship Ship) Collides(otherShip Ship) bool {
	for _, outerCoord := range ship.Coordinates() {
		for _, innerCoord := range otherShip.Coordinates() {
			if innerCoord.x == outerCoord.x && innerCoord.y == outerCoord.y {
				return true
			}
		}
	}
	return false
}

func (ship Ship) Hits() int {
	hits := 0
	for _, s := range ship.structure {
		if !s.healthy {
			hits++
		}
	}
	return hits
}

func (ship Ship) Destroyed() bool {
	return ship.Hits() == ship.class.length
}
