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

type Ship struct {
	class     class
	structure []structureUnit
}

type class struct {
	name   string
	length int
	symbol rune
}

var symbolMap = map[string]rune{
	"Submarine": 'S',
	"Frigate":   'F',
	"Destroyer": 'D',
	"Cruiser":   'C',
	"Carrier":   'T',
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

func NewShip(className string, x int, y int, orientation string) *Ship {
	s := &Ship{class{className, lengthMap[className], symbolMap[className]}, []structureUnit{}}
	o := orientationMap[orientation]
	for i := 0; i < s.Length(); i++ {
		x := x + i*int(o.x)
		y := y + i*int(o.y)
		s.structure = append(s.structure, structureUnit{c: coordinate{x, y}, healthy: true})
	}
	return s
}

func OrientationFromString(orientationString string) (orientation, error) {
	if o, ok := orientationMap[orientationString]; ok {
		return o, nil
	}
	return orientation{0, 0}, fmt.Errorf("cannot convert orientation string: %s", orientationString)
}

func (coordinate coordinate) String() string {
	return fmt.Sprintf("x:%d/y:%d", coordinate.x, coordinate.y)
}

func (coordinate coordinate) X() int {
	return coordinate.x
}

func (coordinate coordinate) Y() int {
	return coordinate.y
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
	return fmt.Sprintf("%v (Stern: %v, Heading: %v, Length: %d, Hits: %d)", ship.class.name, ship.SternCoordinate().String(), ship.Orientation().String(), ship.Length(), ship.Hits())
}

func (ship Ship) Length() int {
	return ship.class.length
}

func (ship Ship) Symbol() rune {
	return ship.class.symbol
}

func (ship Ship) Coordinates() []coordinate {
	var retval = []coordinate{}
	for i := 0; i < ship.Length(); i++ {
		x := ship.SternCoordinate().x + i*int(ship.Orientation().x)
		y := ship.SternCoordinate().y + i*int(ship.Orientation().y)
		retval = append(retval, coordinate{x, y})
	}
	return retval
}

func (ship Ship) BowCoordinate() coordinate {
	return ship.structure[len(ship.structure)-1].c
}

func (ship Ship) SternCoordinate() coordinate {
	return ship.structure[0].c
}

func (ship Ship) Orientation() orientation {
	if ship.Length() < 2 {
		//default to something?
		return orientation{1, 0}
	}
	x := ship.structure[1].c.x - ship.SternCoordinate().x
	y := ship.structure[1].c.y - ship.SternCoordinate().y
	return orientation{int8(x), int8(y)}
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
