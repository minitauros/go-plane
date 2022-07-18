package plane

import (
	"strings"
)

// GetRender is a utility function to render a given surface to stdout.
func GetRender(s *Surface) string {
	rows := make([][]string, 0, s.width)
	for y := s.height - 1; y >= 0; y-- {
		vals := make([]string, 0, s.height)

		for x := 0; x < s.width; x++ {
			if s.IsFilled(Coord{x, y}) {
				vals = append(vals, "x")
			} else {
				vals = append(vals, ".")
			}
		}

		rows = append(rows, vals)
	}

	rowVals := make([]string, 0, len(rows))
	for _, row := range rows {
		rowVals = append(rowVals, strings.Join(row, " "))
	}
	return "\n" + strings.Join(rowVals, "\n")
}
