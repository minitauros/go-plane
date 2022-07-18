package plane

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_FloodFiller_CountUnfilledSurface(t *testing.T) {
	Convey("FloodFiller.Fill()", t, func() {
		s := NewSurface(9, 9)
		filler := NewFloodFiller(s)
		totalSurface := s.width * s.height

		Convey("If counting full board", func() {
			Convey("Without obstacles", func() {
				Convey("Works correctly", func() {
					So(filler.Fill(Coord{0, 0}, Coord{0, 1}), ShouldEqual, totalSurface-1)
				})
			})

			Convey("With obstacles", func() {
				coords := []Coord{
					{1, 1},
					// {2, 1}, // This is the opening into an enclosed space. We expect the inside to be counted.
					{3, 1},
					{4, 1},
					{5, 1},
					{1, 2},
					{1, 3},
					{1, 4},
					{1, 5},
					{2, 5},
					{3, 5},
					{4, 5},
					{5, 5},
					{5, 2},
					{5, 3},
					{5, 4},
				}
				s.Fill(coords...)

				Convey("Works correctly", func() {
					So(filler.Fill(Coord{0, 0}, Coord{0, 1}), ShouldEqual, totalSurface-len(coords)-1)
				})
			})
		})

		Convey("If counting enclosed space within the surface", func() {
			s.Fill(
				Coord{1, 1},
				// Coord{2, 1}, // We start here.
				Coord{3, 1},
				Coord{4, 1},
				Coord{5, 1},
				Coord{1, 2},
				Coord{1, 3},
				Coord{1, 4},
				Coord{1, 5},
				Coord{2, 5},
				Coord{3, 5},
				Coord{4, 5},
				Coord{5, 5},
				Coord{5, 2},
				Coord{5, 3},
				Coord{5, 4},
			)

			start := Coord{2, 1}

			Convey("Without obstacles", func() {
				Convey("Works correctly", func() {
					So(filler.Fill(start, start.GetCoordInDirection(Top)), ShouldEqual, 9)
				})
			})

			Convey("With obstacles", func() {
				s.Fill(
					Coord{2, 4},
					Coord{2, 3},
					Coord{3, 4},
				)

				Convey("Works correctly", func() {
					So(filler.Fill(start, start.GetCoordInDirection(Top)), ShouldEqual, 6)
				})
			})
		})

		Convey("If starting count on an unconnected coordinate, returns 0", func() {
			So(filler.Fill(Coord{0, 0}, Coord{1, 1}), ShouldEqual, 0)
		})

		Convey("If starting on an already filled piece, returns 0", func() {
			s.Fill(Coord{1, 1})

			So(filler.Fill(Coord{0, 0}, Coord{1, 1}), ShouldEqual, 0)
		})
	})
}

func Test_FloodFiller_canReach(t *testing.T) {
	Convey("FloodFiller.canReach()", t, func() {
		s := NewSurface(5, 5)
		filler := NewFloodFiller(s)
		base := Coord{0, 0}
		target := Coord{4, 4}
		possibleStartingCoords := []Coord{{0, 1}, {1, 0}}

		drawHorizontalLineInMiddle := func() {
			for x := 0; x < s.width; x++ {
				s.Fill(Coord{x, s.height / 2})
			}
		}

		Convey("If none of the possible starting coords is filled", func() {
			Convey("If can make path", func() {
				Convey("Returns the number of steps", func() {
					So(filler.canReach(base, target), ShouldBeTrue)
				})
			})

			Convey("If cannot make path", func() {
				drawHorizontalLineInMiddle()

				Convey("Returns false", func() {
					So(filler.canReach(base, target), ShouldBeFalse)
				})
			})
		})

		for _, startingCoord := range possibleStartingCoords {
			Convey(fmt.Sprintf("If one of the possible starting coords is filled (%s)", startingCoord), func() {
				s.Fill(startingCoord)

				Convey("If can make path", func() {
					Convey("Returns the number of steps", func() {
						So(filler.canReach(base, target), ShouldBeTrue)
					})
				})

				Convey("If cannot make path", func() {
					drawHorizontalLineInMiddle()

					Convey("Returns false", func() {
						So(filler.canReach(base, target), ShouldBeFalse)
					})
				})
			})
		}

		Convey("If both the possible starting coords are filled", func() {
			s.Fill(possibleStartingCoords...)

			Convey("Returns false", func() {
				So(filler.canReach(base, target), ShouldBeFalse)
			})
		})

		Convey("With obstacles, but the path not obstructed fully", func() {
			s.Fill(Coord{2, 2}, Coord{3, 4})

			Convey("Returns the number of steps", func() {
				So(filler.canReach(base, target), ShouldBeTrue)
			})
		})

		Convey("If allowed to start only in one direction", func() {
			mayStartOnlyAt := Coord{1, 0}

			// (S = start, T = target, x = filled)
			// Make area look like
			// . . . . T
			// . . . x .
			// . . . x .
			// . x x x .
			// S . . . .
			// It should be possible to reach T by going to the top, but we don't allow it.
			for x := 1; x < s.width-1; x++ {
				s.Fill(Coord{x, 1})
			}
			for y := 1; y < s.height-1; y++ {
				s.Fill(Coord{3, y})
			}

			Convey("If can make path", func() {
				Convey("Returns the number of steps", func() {
					So(filler.canReach(base, target, mayStartOnlyAt), ShouldBeTrue)
				})
			})

			Convey("If cannot make path", func() {
				// fill bottom right coord.
				s.Fill(Coord{s.width - 1, 0})

				Convey("Returns false", func() {
					So(filler.canReach(base, target, mayStartOnlyAt), ShouldBeFalse)
				})

				Convey("To the top still works (where we can make a path)", func() {
					So(filler.canReach(base, target, Coord{0, 1}), ShouldBeTrue)
				})
			})
		})
	})
}

func Test_FloodFiller_CountSteps(t *testing.T) {
	Convey("FloodFiller.CountSteps()", t, func() {
		s := NewSurface(5, 5)
		filler := NewFloodFiller(s)
		base := Coord{0, 0}
		target := Coord{4, 4}
		possibleStartingCoords := []Coord{{0, 1}, {1, 0}}

		drawHorizontalLineInMiddle := func() {
			for x := 0; x < s.width; x++ {
				s.Fill(Coord{x, s.height / 2})
			}
		}

		for _, startingCoord := range possibleStartingCoords {
			Convey(fmt.Sprintf("If one of the possible starting coords is filled (%s)", startingCoord), func() {
				s.Fill(startingCoord)

				Convey("If can make path", func() {
					Convey("Returns the shortest path", func() {
						So(filler.CountSteps(base, target), ShouldEqual, s.width+s.height-1)
					})
				})

				Convey("If cannot make path", func() {
					drawHorizontalLineInMiddle()

					Convey("Returns -1", func() {
						So(filler.CountSteps(base, target), ShouldEqual, -1)
					})
				})
			})
		}

		Convey("When target is around an obstacle", func() {
			// (S = start, T = target, x = filled)
			// Make area look like
			// T . . . . . .
			// . . . x x x .
			// x x . . . x .
			// . x x x . x .
			// S . . . . . .
			// We expect the filler to find the shortest route (12).
			rows := [][]int{
				{0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 1, 1, 1, 0},
				{1, 1, 0, 0, 0, 1, 0},
				{0, 1, 1, 1, 0, 1, 0},
				{0, 0, 0, 0, 0, 0, 0},
			}
			s := NewSurface(7, 5)
			for y, row := range rows {
				for x, val := range row {
					if val == 0 {
						continue
					}
					s.Fill(Coord{x, s.height - y - 1})
				}
			}
			filler := NewFloodFiller(s)
			base := Coord{0, 0}
			target := Coord{0, 5}

			Convey("If can make path", func() {
				Convey("Returns the number of steps", func() {
					So(filler.CountSteps(base, target), ShouldEqual, -1)
				})
			})

			Convey("If cannot make path", func() {
				s.Fill(Coord{0, 1})

				Convey("Returns false", func() {
					So(filler.CountSteps(base, target), ShouldEqual, -1)
				})
			})
		})
	})
}
