# Plane

This is a package that was created during a "hackday" at which we attempted to create a [Battlesnake](https://play.battlesnake.com/) and win a battle against each other. The Battlesnake [documentation](https://docs.battlesnake.com/references/useful-algorithms) recommended to use a [flood fill algorithm](https://en.wikipedia.org/wiki/Flood_fill), and I couldn't find an implementation in Go that I could use.

This package contains methods for working with a surface (to be able to tell which coordinates on a surface are/aren't filled) and for flood filling such surfaces, to determine for example how much unfilled surface exists within a possibly enclosed space, or what the quickest path is - even around obstacles - from A to B.


## Surface

```go
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

	// Looping over all filled coords.
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

```

## Flood filler

```go
package main

import (
	"github.com/minitauros/go-plane"
)

func main() {
	// Flood filler needs a surface to work with.
	surface := plane.NewSurface(5, 5)

	// Create new flood filler.
	// We clone the surface, because flood filler will fill the surface,
	// and we do not want to mess up the original state.
	ff := plane.NewFloodFiller(surface.Clone())

	// Fill the plane, using 0,0 as base and starting the flood at 0,1.
	// Note that this does **not** fill `base`.
	// That means that the whole surface will be filled after this fill,
	// except the `base` coordinate. If you want to fill this,
	// call `surface.Fill()`.
	filledCoords := ff.Flood(plane.Coord{0, 0}, plane.Coord{0, 1})

	// Return the quickest path from 0,0 to 4,4.
	// This will go around obstacles.
	ff.CountSteps(plane.Coord{0, 0}, plane.Coord{4, 4}) // 9

	// Returns true if 5,5 can be reached, i.e. if there are no
	// obstacles (filled coords) in the way that the flood cannot pass
	// through in some way.
	ff.CanReach(plane.Coord{0, 0}, plane.Coord{4, 4}) // True

	// Returns true if 5,5 can be reached, i.e. if there are no
	// obstacles (filled coords) in the way that the flood cannot pass
	// through in some way.
	// Forces the flood to start at 0,1 and gives it no other options.
	ff.CanReachWhenStartingFloodAt(plane.Coord{0, 0}, plane.Coord{4, 4}, plane.Coord{0, 1})
}

```

## Notes

This package was created while working under time pressure, because I had to win the Battlesnake hackathon. I have added tests for some cases, but some are missing. So far code seems to be working. 