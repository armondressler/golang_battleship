package weapon

type coordinate struct {
	x, y int
}

type Exploder interface {
	Explode(coordinate) []coordinate
	Symbol() rune
}

type weapon struct {
	name   string
	symbol rune
}

type simpleTorpedo struct {
	weapon
}

type seaMine struct {
	weapon
}

func (weapon weapon) Symbol() rune {
	return weapon.symbol
}

func (s seaMine) Explode(c coordinate) []coordinate {
	affectedCoordinates := []coordinate{}
	affectedCoordinates = append(affectedCoordinates, c)
	return affectedCoordinates
}

func (t simpleTorpedo) Explode(c coordinate) []coordinate {
	affectedCoordinates := []coordinate{}
	affectedCoordinates = append(affectedCoordinates, c)
	return affectedCoordinates
}
