package plane

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_FloodFiller_Fill(t *testing.T) {
	Convey("FloodFiller.Flood()", t, func() {
		s := NewSurface(9, 9)
		filler := NewFloodFiller(s)

		Convey("If counting full board", func() {
			Convey("Without obstacles", func() {
				Convey("Works correctly", func() {
					filled := filler.Flood(Coord{0, 0}, Coord{0, 1})
					for x := 0; x < s.width; x++ {
						for y := 0; y < s.height; y++ {
							if x == 0 && y == 0 {
								// We don't expect the starting coord to be filled.
								continue
							}
							var isFilled bool
							for _, coord := range filled {
								if coord.X == x && coord.Y == y {
									isFilled = true
									break
								}
							}
							Convey(fmt.Sprintf("%d,%d is filled", x, y), func() {
								So(isFilled, ShouldBeTrue)
							})
						}
					}
					So(filled, ShouldHaveLength, s.width*s.height-1)
				})
			})

			Convey("With obstacles", func() {
				// 08 | . . . . . . . . .
				// 07 | . . . . . . . . .
				// 06 | . . . . . . . . .
				// 05 | . x x x x x . . .
				// 04 | . x . . . x . . .
				// 03 | . x . . . x . . .
				// 02 | . x . . . x . . .
				// 01 | . x . x x x . . .
				// 00 | . . . . . . . . .
				//     ------------------
				//      0 1 2 3 4 5 6 7 8
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
					filled := filler.Flood(Coord{0, 0}, Coord{0, 1})

					expectedFilled := Coords{
						// Enclosed space + opening
						{2, 1}, {2, 2}, {3, 2}, {4, 2},
						{2, 3}, {3, 3}, {4, 3},
						{2, 4}, {3, 4}, {4, 4},
						// Outside area.
						{1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}, {6, 0}, {7, 0}, {8, 0},
						{0, 1}, {6, 1}, {7, 1}, {8, 1},
						{0, 2}, {6, 2}, {7, 2}, {8, 2},
						{0, 3}, {6, 3}, {7, 3}, {8, 3},
						{0, 4}, {6, 4}, {7, 4}, {8, 4},
						{0, 5}, {6, 5}, {7, 5}, {8, 5},
						{0, 6}, {1, 6}, {2, 6}, {3, 6}, {4, 6}, {5, 6}, {6, 6}, {7, 6}, {8, 6},
						{0, 7}, {1, 7}, {2, 7}, {3, 7}, {4, 7}, {5, 7}, {6, 7}, {7, 7}, {8, 7},
						{0, 8}, {1, 8}, {2, 8}, {3, 8}, {4, 8}, {5, 8}, {6, 8}, {7, 8}, {8, 8},
					}

					for _, coord := range expectedFilled {
						var isFilled bool
						for _, coord2 := range filled {
							if coord.Equals(coord2) {
								isFilled = true
								break
							}
						}
						Convey(fmt.Sprintf("Expected %s is filled", coord), func() {
							So(isFilled, ShouldBeTrue)
						})
					}

					for _, coord := range filled {
						var isExpected bool
						for _, coord2 := range expectedFilled {
							if coord.Equals(coord2) {
								isExpected = true
								break
							}
						}
						Convey(fmt.Sprintf("Flood of %s was expected", coord), func() {
							So(isExpected, ShouldBeTrue)
						})
					}
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
					So(len(filler.Flood(start, start.GetCoordInDirection(Top))), ShouldEqual, 9)
				})
			})

			Convey("With obstacles", func() {
				s.Fill(
					Coord{2, 4},
					Coord{2, 3},
					Coord{3, 4},
				)

				Convey("Works correctly", func() {
					So(len(filler.Flood(start, start.GetCoordInDirection(Top))), ShouldEqual, 6)
				})
			})
		})

		Convey("If starting count on an unconnected coordinate, returns 0", func() {
			So(len(filler.Flood(Coord{0, 0}, Coord{1, 1})), ShouldEqual, 0)
		})

		Convey("If starting on an already filled piece, returns 0", func() {
			s.Fill(Coord{1, 1})

			So(len(filler.Flood(Coord{0, 0}, Coord{1, 1})), ShouldEqual, 0)
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
				Convey("Returns true", func() {
					So(filler.canReach(base, target, false), ShouldBeTrue)
				})
			})

			Convey("If cannot make path", func() {
				drawHorizontalLineInMiddle()

				Convey("Returns false", func() {
					So(filler.canReach(base, target, false), ShouldBeFalse)
				})
			})
		})

		for _, startingCoord := range possibleStartingCoords {
			Convey(fmt.Sprintf("If one of the possible starting coords is filled (%s)", startingCoord), func() {
				s.Fill(startingCoord)

				Convey("If can make path", func() {
					Convey("Returns true", func() {
						So(filler.canReach(base, target, false), ShouldBeTrue)
					})
				})

				Convey("If cannot make path", func() {
					drawHorizontalLineInMiddle()

					Convey("Returns false", func() {
						So(filler.canReach(base, target, false), ShouldBeFalse)
					})
				})
			})
		}

		Convey("If both the possible starting coords are filled", func() {
			s.Fill(possibleStartingCoords...)

			Convey("Returns false", func() {
				So(filler.canReach(base, target, false), ShouldBeFalse)
			})
		})

		Convey("With obstacles, but the path not obstructed fully", func() {
			s.Fill(Coord{2, 2}, Coord{3, 4})

			Convey("Returns the number of steps", func() {
				So(filler.canReach(base, target, false), ShouldBeTrue)
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
				Convey("Returns true", func() {
					So(filler.canReach(base, target, false, mayStartOnlyAt), ShouldBeTrue)
				})
			})

			Convey("If cannot make path", func() {
				// flood bottom right coord.
				s.Fill(Coord{s.width - 1, 0})

				Convey("Returns false", func() {
					So(filler.canReach(base, target, false, mayStartOnlyAt), ShouldBeFalse)
				})

				Convey("To the top still works (where we can make a path)", func() {
					So(filler.canReach(base, target, false, Coord{0, 1}), ShouldBeTrue)
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
						So(filler.CountSteps(base, target), ShouldEqual, s.width+s.height-2)
					})
				})

				Convey("If cannot make path", func() {
					drawHorizontalLineInMiddle()

					Convey("Returns -1", func() {
						steps := filler.CountSteps(base, target)

						So(steps, ShouldEqual, -1)
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
			s := NewSurface(7, 5)
			s.fillRows([][]int{
				{0, 0, 0, 0, 0, 0, 0},
				{0, 0, 0, 1, 1, 1, 0},
				{1, 1, 0, 0, 0, 1, 0},
				{0, 1, 1, 1, 0, 1, 0},
				{0, 0, 0, 0, 0, 0, 0},
			})
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

		Convey("If target is next to starting point", func() {
			// . . .
			// . . .
			// s t .
			// s = start; t = target
			s := NewSurface(3, 3)
			filler := NewFloodFiller(s)

			Convey("And is not filled already", func() {
				Convey("Returns 1", func() {
					numSteps := filler.CountSteps(Coord{0, 0}, Coord{1, 0})

					So(numSteps, ShouldEqual, 1)
				})
			})

			Convey("And is filled already", func() {
				s.Fill(Coord{1, 0})

				Convey("Returns 1", func() {
					numSteps := filler.CountSteps(Coord{0, 0}, Coord{1, 0})

					So(numSteps, ShouldEqual, 1)
				})

			})
		})
	})
}

func Test_FloodFiller_flood(t *testing.T) {
	Convey("FloodFiller.flood()", t, func() {
		s := NewSurface(3, 3)
		filler := NewFloodFiller(s)

		Convey("When tracking distance", func() {
			Convey("Correctly fills the field", func() {
				filled := filler.flood(Coord{0, 0}, Coord{1, 0}, true)

				So(filled, ShouldResemble, Coords{})

				// Bottom row
				v, ok := s.surface[0][0]
				So(ok, ShouldBeFalse)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 0,
				})

				v, ok = s.surface[1][0]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 1,
				})

				v, ok = s.surface[2][0]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 2,
				})

				// Middle row
				v, ok = s.surface[0][1]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 3,
				})

				v, ok = s.surface[1][1]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 2,
				})

				v, ok = s.surface[2][1]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 3,
				})

				// Top row
				v, ok = s.surface[0][2]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 4,
				})

				v, ok = s.surface[1][2]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 3,
				})

				v, ok = s.surface[2][2]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 4,
				})
			})

			Convey("Correctly fills already filled coords (that don't have a value)", func() {
				s.Fill(Coord{0, 1})

				filler.flood(Coord{0, 0}, Coord{1, 0}, true)

				// Bottom row
				v, ok := s.surface[0][0]
				So(ok, ShouldBeFalse)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 0,
				})

				v, ok = s.surface[1][0]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 1,
				})

				v, ok = s.surface[2][0]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 2,
				})

				// Middle row
				v, ok = s.surface[0][1]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: true,
					distance: 3,
				})

				v, ok = s.surface[1][1]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 2,
				})

				v, ok = s.surface[2][1]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 3,
				})

				// Top row
				v, ok = s.surface[0][2]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 4,
				})

				v, ok = s.surface[1][2]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 3,
				})

				v, ok = s.surface[2][2]
				So(ok, ShouldBeTrue)
				So(v, ShouldResemble, coordVal{
					isFilled: false,
					distance: 4,
				})
			})
		})
	})
}
