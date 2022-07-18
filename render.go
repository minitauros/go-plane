package plane

import (
	"fmt"
	"strconv"
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
	for i, row := range rows {
		row = append([]string{fmt.Sprintf("%01d | ", s.height-i-1)}, row...)
		rowVals = append(rowVals, strings.Join(row, " "))
	}

	rowVals = append(rowVals, "     "+strings.Repeat("-", s.width*2-1))

	xLegendVals := []string{"    "}
	for x := 0; x < s.width; x++ {
		xLegendVals = append(xLegendVals, strconv.Itoa(x))
	}
	rowVals = append(rowVals, strings.Join(xLegendVals, " "))

	return "\n" + strings.Join(rowVals, "\n")
}
