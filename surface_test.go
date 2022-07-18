package plane

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSurface(t *testing.T) {
	type input struct {
		width  int
		height int
	}
	testCases := []struct {
		description string
		input       input
		expected    *Surface
	}{
		{
			description: "",
			input: input{
				width:  3,
				height: 3,
			},
			expected: &Surface{
				width:  3,
				height: 3,
				surface: surfaceMap{
					0: {},
					1: {},
					2: {},
				},
			},
		},
	}

	Convey("NewSurface()", t, func() {
		for i, tc := range testCases {
			Convey(fmt.Sprintf("%d: %s", i, tc.description), func() {
				res := NewSurface(tc.input.width, tc.input.height)

				So(res.width, ShouldEqual, tc.expected.width)
				So(res.height, ShouldEqual, tc.expected.height)
				So(res.surface, ShouldResemble, tc.expected.surface)
				So(len(res.surface), ShouldEqual, tc.input.width)
			})
		}
	})
}

func TestSurface_fill(t *testing.T) {
	// Assuming a surface of 3x3.
	testCases := []struct {
		description string
		input       []Coord
		expected    surfaceMap
	}{
		{
			description: "sets all corners correctly",
			input:       []Coord{{0, 0}, {0, 2}, {2, 2}, {2, 0}},
			expected: surfaceMap{
				0: {
					0: 0,
					2: 0,
				},
				1: {},
				2: {
					0: 0,
					2: 0,
				},
			},
		},
		{
			description: "somewhere in the middle",
			input:       []Coord{{1, 1}, {1, 2}, {2, 2}},
			expected: surfaceMap{
				0: {},
				1: {
					1: 0,
					2: 0,
				},
				2: {
					2: 0,
				},
			},
		},
		{
			description: "does not set coords that are out of bounds",
			input:       []Coord{{-1, -1}, {4, 4}},
			expected: surfaceMap{
				0: {},
				1: {},
				2: {},
			},
		},
	}

	Convey("fill()", t, func() {
		for i, tc := range testCases {
			Convey(fmt.Sprintf("%d: %s", i, tc.description), func() {
				s := NewSurface(3, 3)
				s.Fill(tc.input...)

				for x := 0; x < s.width; x++ {
					for y := 0; y < s.height; y++ {
						var hasCoord bool
						for _, coord := range tc.input {
							if coord.X == x && coord.Y == y {
								hasCoord = true
								break
							}
						}

						_, ok := s.surface[x][y]
						if hasCoord {
							So(ok, ShouldBeTrue)
						} else {
							So(ok, ShouldBeFalse)
						}
					}
				}
			})
		}
	})
}

func Test_StandardSurface_Remove(t *testing.T) {
	type input struct {
		surface surfaceMap
		remove  []Coord
	}
	testCases := []struct {
		description string
		input       input
		expected    surfaceMap
	}{
		{
			description: "removes correctly",
			input: input{
				surface: surfaceMap{
					1: map[int]int{
						2: 0,
					},
				},
				remove: []Coord{{1, 2}},
			},
			expected: surfaceMap{
				1: map[int]int{},
			},
		},
		{
			description: "does nothing when given coord to remove does not exist",
			input: input{
				surface: surfaceMap{
					4: map[int]int{
						9: 0,
					},
				},
				remove: []Coord{{1, 2}},
			},
			expected: surfaceMap{
				4: map[int]int{
					9: 0,
				},
			},
		},
		{
			description: "removes from the middle correctly",
			input: input{
				surface: surfaceMap{
					1: map[int]int{
						2: 0,
					},
					4: map[int]int{
						5: 0,
						6: 0,
						9: 0,
					},
					8: map[int]int{
						9: 0,
					},
				},
				remove: []Coord{{4, 6}},
			},
			expected: surfaceMap{
				1: map[int]int{
					2: 0,
				},
				4: map[int]int{
					5: 0,
					9: 0,
				},
				8: map[int]int{
					9: 0,
				},
			},
		},
		{
			description: "removes multiple correctly",
			input: input{
				surface: surfaceMap{
					1: map[int]int{
						2: 0,
					},
					4: map[int]int{
						5: 0,
						6: 0,
						9: 0,
					},
					8: map[int]int{
						9: 0,
					},
				},
				remove: []Coord{{4, 6}, {8, 9}},
			},
			expected: surfaceMap{
				1: map[int]int{
					2: 0,
				},
				4: map[int]int{
					5: 0,
					9: 0,
				},
				8: map[int]int{},
			},
		},
	}

	Convey("Surface.Remove()", t, func() {
		for i, tc := range testCases {
			Convey(fmt.Sprintf("%d: %s", i, tc.description), func() {
				s := NewSurface(0, 0)
				s.surface = tc.input.surface
				s.Remove(tc.input.remove...)

				So(tc.input.surface, ShouldResemble, tc.expected)
			})
		}
	})
}
