package plane

import (
	"math"
	"strconv"
	"strings"
)

// Coord is a coordinate.
type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Equals returns true if the current coord equals the given other coord.
func (c Coord) Equals(other Coord) bool {
	return c.X == other.X && c.Y == other.Y
}

// String satisfies stringer.
func (c Coord) String() string {
	return strconv.Itoa(c.X) + "," + strconv.Itoa(c.Y)
}

// ConnectsTo returns true if the current coord connects directly to the given other coordinate.
func (c Coord) ConnectsTo(other Coord) bool {
	horizontalDiff := abs(c.X - other.X)
	verticalDiff := abs(c.Y - other.Y)
	return horizontalDiff <= 1 && verticalDiff <= 1 && !(horizontalDiff == 1 && verticalDiff == 1)
}

// GetCoordInDirection returns the first coordinate in the given direction from the current coordinate.
func (c Coord) GetCoordInDirection(d Direction) Coord {
	switch d {
	case Top:
		return Coord{c.X, c.Y + 1}
	case Bot:
		return Coord{c.X, c.Y - 1}
	case Right:
		return Coord{c.X + 1, c.Y}
	case Left:
		return Coord{c.X - 1, c.Y}
	}
	return c
}

// GetCoordsTo returns the coords of which the direction can be taken to move towards the given `to` coord.
func (c Coord) GetCoordsTo(to Coord) Coords {
	coords := make(Coords, 0, 2)
	if to.X > c.X {
		coords = append(coords, c.GetCoordInDirection(Right))
	} else if to.X < c.X {
		coords = append(coords, c.GetCoordInDirection(Left))
	}
	if to.Y > c.Y {
		coords = append(coords, c.GetCoordInDirection(Top))
	} else if to.Y < c.Y {
		coords = append(coords, c.GetCoordInDirection(Bot))
	}
	return coords
}

// GetDirectionsTo returns the directions that can be taken to move towards the given `to` coord.
func (c Coord) GetDirectionsTo(to Coord) []Direction {
	dirs := make([]Direction, 0, 2)
	if to.X > c.X {
		dirs = append(dirs, Right)
	} else if to.X < c.X {
		dirs = append(dirs, Left)
	}
	if to.Y > c.Y {
		dirs = append(dirs, Top)
	} else if to.Y < c.Y {
		dirs = append(dirs, Bot)
	}
	return dirs
}

// GetCoordsAround returns all coords around the current coord.
func (c Coord) GetCoordsAround() Coords {
	return []Coord{
		{c.X + 1, c.Y},
		{c.X - 1, c.Y},
		{c.X, c.Y + 1},
		{c.X, c.Y - 1},
	}
}

// GetCoordAt returns the coord at the given offset from the current coord.
func (c Coord) GetCoordAt(xOffset, yOffset int) Coord {
	return Coord{c.X + xOffset, c.Y + yOffset}
}

// Coords is an array of coordinates.
type Coords []Coord

// GetIntersections returns of the given coordinates those that intersect with the current ones.
func (coords Coords) GetIntersections(other []Coord) Coords {
	var intersections []Coord
	for _, coord1 := range coords {
		for _, coord2 := range other {
			if coord1.Equals(coord2) {
				intersections = append(intersections, coord1)
			}
		}
	}
	return intersections
}

// Remove removes the given coords from the current ones.
func (coords *Coords) Remove(coordsToRemove ...Coord) {
	if len(*coords) == 0 || len(coordsToRemove) == 0 {
		return
	}
	var newCoords Coords
	for _, c := range *coords {
		var mustRemove bool
		for _, cToRemove := range coordsToRemove {
			if c.Equals(cToRemove) {
				mustRemove = true
			}
		}
		if !mustRemove {
			newCoords = append(newCoords, c)
		}
	}
	*coords = newCoords
}

// Contains returns true if the given coord is in the current collection of coords.
func (coords Coords) Contains(coord Coord) bool {
	for _, c := range coords {
		if c.Equals(coord) {
			return true
		}
	}
	return false
}

// String satisfies stringer.
func (coords Coords) String() string {
	chunks := make([]string, 0, len(coords))
	for _, coord := range coords {
		chunks = append(chunks, coord.String())
	}
	return strings.Join(chunks, " ")
}

func abs(in int) int {
	return int(math.Abs(float64(in)))
}
