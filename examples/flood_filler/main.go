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

	// Flood the plane, using 0,0 as base and starting the flood at 0,1.
	// Note that this does **not** fill `base`.
	// That means that the whole surface will be filled after this fill,
	// except the `base` coordinate. If you want to fill this,
	// call `surface.Flood()`.
	filledCoords := ff.Flood(plane.Coord{0, 0}, plane.Coord{0, 1})

	// Return the quickest path from 0,0 to 4,4.
	// This will go around obstacles.
	ff.CountSteps(plane.Coord{0, 0}, plane.Coord{0, 1}) // 9

	// Returns true if 5,5 can be reached, i.e. if there are no
	// obstacles (filled coords) in the way that the flood cannot pass
	// through in some way.
	ff.CanReach(plane.Coord{0, 0}, plane.Coord{4, 4}) // True

	// Returns true if 5,5 can be reached, i.e. if there are no
	// obstacles (filled coords) in the way that the flood cannot pass
	// through in some way.
	// Forces the flood to start at 0,1 and gives it no other options.
	ff.CanReachWhenStartingFillAt(plane.Coord{0, 0}, plane.Coord{4, 4}, plane.Coord{0, 1})
}
