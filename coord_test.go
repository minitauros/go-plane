package plane

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Coords_GetIntersections(t *testing.T) {
	type input struct {
		base  Coords
		other Coords
	}
	testCases := []struct {
		description string
		input       input
		expected    Coords
	}{
		{
			description: "All the same finds all",
			input: input{
				base:  Coords{{1, 1}, {1, 2}, {3, 4}, {3, 1}},
				other: Coords{{1, 1}, {1, 2}, {3, 4}, {3, 1}},
			},
			expected: Coords{{1, 1}, {1, 2}, {3, 4}, {3, 1}},
		},
		{
			description: "First and last the same finds correctly",
			input: input{
				base:  Coords{{1, 1}, {1, 2}, {3, 4}, {3, 1}},
				other: Coords{{1, 1}, {3, 1}},
			},
			expected: Coords{{1, 1}, {3, 1}},
		},
		{
			description: "Only last the same finds correctly",
			input: input{
				base:  Coords{{3, 1}},
				other: Coords{{1, 1}, {1, 2}, {3, 4}, {3, 1}},
			},
			expected: Coords{{3, 1}},
		},
	}

	Convey("Coords_GetIntersections()", t, func() {
		for i, tc := range testCases {
			Convey(fmt.Sprintf("%d: %s", i, tc.description), func() {
				intersections := tc.input.base.GetIntersections(tc.input.other)

				So(intersections, ShouldResemble, tc.expected)
			})
		}
	})
}

func Test_Coords_Remove(t *testing.T) {
	type input struct {
		coords   Coords
		toRemove Coords
	}
	testCases := []struct {
		description string
		input       input
		expected    Coords
	}{
		{
			description: "remove from end",
			input: input{
				coords:   Coords{{1, 1}, {1, 2}, {2, 1}, {0, 0}},
				toRemove: Coords{{0, 0}},
			},
			expected: Coords{{1, 1}, {1, 2}, {2, 1}},
		},
		{
			description: "remove from start",
			input: input{
				coords:   Coords{{1, 1}, {1, 2}, {2, 1}, {0, 0}},
				toRemove: Coords{{1, 1}},
			},
			expected: Coords{{1, 2}, {2, 1}, {0, 0}},
		},
		{
			description: "remove from middle",
			input: input{
				coords:   Coords{{1, 1}, {1, 2}, {2, 1}, {0, 0}},
				toRemove: Coords{{1, 2}},
			},
			expected: Coords{{1, 1}, {2, 1}, {0, 0}},
		},
		{
			description: "remove multiple from middle",
			input: input{
				coords:   Coords{{1, 1}, {1, 2}, {2, 1}, {0, 0}},
				toRemove: Coords{{1, 2}, {2, 1}},
			},
			expected: Coords{{1, 1}, {0, 0}},
		},
		{
			description: "remove nothing",
			input: input{
				coords:   Coords{{1, 1}, {1, 2}, {2, 1}, {0, 0}},
				toRemove: Coords{{1, 5}},
			},
			expected: Coords{{1, 1}, {1, 2}, {2, 1}, {0, 0}},
		},
	}

	Convey("Coords.Remove()", t, func() {
		for i, tc := range testCases {
			Convey(fmt.Sprintf("%d: %s", i, tc.description), func() {
				tc.input.coords.Remove(tc.input.toRemove...)

				So(tc.input.coords, ShouldResemble, tc.expected)
			})
		}
	})
}

func Test_Coord_ConnectsTo(t *testing.T) {
	type input struct {
		base  Coord
		other Coord
	}
	testCases := []struct {
		description string
		input       input
		expected    bool
	}{
		{
			description: "if connected to the top, returns true",
			input: input{
				base:  Coord{1, 1},
				other: Coord{1, 2},
			},
			expected: true,
		},
		{
			description: "if connected to the right, returns true",
			input: input{
				base:  Coord{1, 1},
				other: Coord{2, 1},
			},
			expected: true,
		},
		{
			description: "if connected to the bot, returns true",
			input: input{
				base:  Coord{1, 1},
				other: Coord{1, 0},
			},
			expected: true,
		},
		{
			description: "if connected to the left, returns true",
			input: input{
				base:  Coord{1, 1},
				other: Coord{0, 1},
			},
			expected: true,
		},
		{
			description: "if connected diagonally, returns false",
			input: input{
				base:  Coord{1, 1},
				other: Coord{0, 0},
			},
			expected: false,
		},
		{
			description: "if not connected, returns false",
			input: input{
				base:  Coord{1, 1},
				other: Coord{3, 3},
			},
			expected: false,
		},
	}

	Convey("Coord_ConnectsTo()", t, func() {
		for i, tc := range testCases {
			Convey(fmt.Sprintf("%d: %s", i, tc.description), func() {
				res := tc.input.base.ConnectsTo(tc.input.other)

				So(res, ShouldEqual, tc.expected)
			})
		}
	})
}
