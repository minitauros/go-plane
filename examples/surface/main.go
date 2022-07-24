package main

import (
	"github.com/minitauros/go-plane"
)

func main() {
	// Create new surface.
	// Both the x and y coordinates in this example are in range (0,4).
	surface := plane.NewSurface(5, 5)

	// Creates a fill as follows.
	// 4 | . . . . .
	// 3 | . . . . .
	// 2 | . . . . .
	// 1 | x . . . .
	// 0 | x x . . .
	//   -----------
	//     0 1 2 3 4
	surface.Fill(
		plane.Coord{0, 0},
		plane.Coord{1, 0},
		plane.Coord{0, 1},
	)

	// Checking if a coord is filled.
	surface.IsFilled(plane.Coord{0, 0}) // True

	// Removing/unfilling coords.
	surface.Remove(plane.Coord{0, 0})
	surface.IsFilled(plane.Coord{0, 0}) // False

	// Iterate over all filled coords.
	for coord := range surface.EachFilled() {
		// Do something..
	}

	// Getting all filled coords at once.
	surface.GetFilled() // Coords{{0, 0}, {1, 0}, {0, 1}}

	// Count filled.
	surface.CountFilled() // 3

	// Count unfilled.
	surface.CountUnfilled() // 22

	// Get the total surface.
	surface.TotalSurface() // 25

	// Getting the center coord.
	surface.GetCenter() // plane.Coord{2, 2}

	// Checking if a coord fits.
	surface.Fits(plane.Coord{-1, -1}) // False
	surface.Fits(plane.Coord{0, 0})   // True

	// Clone the surface.
	// This is useful when passing it to the flood filler, as the flood
	// filler will change the surface's state, and you may want to remember
	// the original state.
	surface.Clone() // *Surface
}
