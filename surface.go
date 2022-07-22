package plane

// coordVal describes the values that exist at a coordinate. A coordinate can be filled but have no distance, or have a
// distance but not be filled. This data is saved for different purposes. We want to know if a position is filled when
// we want to flood flood (don't flood the same coord twice). We want to know the distance when calculating the distance.
// In that case, sometimes we have to visit positions twice (because a longer path may have reached a coord before a
// shorter path).
type coordVal struct {
	isFilled bool
	distance int
}

// surfaceMap defines which coordinates in a surface are filled.
// The structure is [x][y]coordVal.
type surfaceMap map[int]map[int]coordVal

// Surface represents a surface of a given width and height.
// For width 5 and height 5, the coordinates would range from 0-4x and 0-4y.
// 0,0 is bottom Left.
type Surface struct {
	width  int
	height int
	// surface keeps track of which coordinates are filled.
	surface surfaceMap
}

// NewSurface returns a new surface.
func NewSurface(width int, height int) *Surface {
	s := make(surfaceMap, width)
	for i := 0; i < width; i++ {
		s[i] = make(map[int]coordVal, height)
	}
	return &Surface{
		width:   width,
		height:  height,
		surface: s,
	}
}

// GetCenter returns the center of the surface. It will round down if the center coordinate is not a round number.
func (s *Surface) GetCenter() Coord {
	return Coord{s.width / 2, s.height / 2}
}

// EachFilled returns each filled coord.
func (s *Surface) EachFilled() <-chan Coord {
	ch := make(chan Coord)
	go func() {
		for x := 0; x < s.width; x++ {
			for y := 0; y < s.height; y++ {
				_, ok := s.surface[x][y]
				if ok {
					ch <- Coord{x, y}
				}
			}
		}
		close(ch)
	}()
	return ch
}

// GetFilled returns all filled coords.
func (s *Surface) GetFilled() Coords {
	var filled Coords
	for coord := range s.EachFilled() {
		filled = append(filled, coord)
	}
	return filled
}

// CountFilled returns the number of filled coords.
func (s *Surface) CountFilled() int {
	var num int
	for range s.EachFilled() {
		num++
	}
	return num
}

// CountUnfilled returns the number of unfilled coords.
func (s *Surface) CountUnfilled() int {
	return s.TotalSurface() - s.CountFilled()
}

// TotalSurface returns the total surface area.
func (s *Surface) TotalSurface() int {
	return s.width * s.height
}

// Fits returns true if the given coord fits on the surface.
func (s *Surface) Fits(coord Coord) bool {
	return coord.X >= 0 && coord.Y >= 0 && coord.X < s.width && coord.Y < s.height
}

// Remove removes (unfills) the given coords from the surface.
func (s *Surface) Remove(coords ...Coord) {
	for _, coord := range coords {
		delete(s.surface[coord.X], coord.Y)
	}
}

// Fill fills the given coords.
func (s *Surface) Fill(coords ...Coord) {
	for _, coord := range coords {
		// If already filled, don't flood again.
		// We don't want to overwrite the value.
		if s.IsFilled(coord) {
			continue
		}
		s.surface[coord.X][coord.Y] = coordVal{
			isFilled: true,
		}
	}
}

// IsFilled returns true if the given coord is filled or does not fit on the surface.
func (s *Surface) IsFilled(coord Coord) bool {
	if !s.Fits(coord) {
		return true
	}
	v, ok := s.surface[coord.X][coord.Y]
	return ok && v.isFilled
}

// Clone returns a clone of the surface.
func (s *Surface) Clone() *Surface {
	clone := NewSurface(s.width, s.height)
	for coord := range s.EachFilled() {
		clone.forceFill(coord)
	}
	return clone
}

func (s *Surface) hasDistance(coord Coord) bool {
	v, ok := s.surface[coord.X][coord.Y]
	return ok && v.distance > 0
}

func (s *Surface) getDistance(coord Coord) int {
	if v, ok := s.surface[coord.X][coord.Y]; ok {
		return v.distance
	}
	return -1
}

func (s *Surface) getValue(coord Coord) (coordVal, bool) {
	v, ok := s.surface[coord.X][coord.Y]
	return v, ok
}

func (s *Surface) setDistance(coord Coord, distance int) {
	v := s.surface[coord.X][coord.Y]
	v.distance = distance
	s.surface[coord.X][coord.Y] = v
}

// forceFill does a regular flood, but does not check if the coords actually fit on the surface,
// which is slightly faster.
func (s *Surface) forceFill(coords ...Coord) {
	for _, coord := range coords {
		v := s.surface[coord.X][coord.Y]
		v.isFilled = true
		s.surface[coord.X][coord.Y] = v
	}
}

func (s *Surface) getCoordsFilledAround(coord Coord) Coords {
	coordsAround := coord.GetCoordsAround()
	filled := make(Coords, 0, 4)
	for _, c := range coordsAround {
		if s.IsFilled(c) {
			filled = append(filled, c)
		}
	}
	return filled
}

// fillRows is a utility method to conveniently flood the surface.
// Example input:
// [][]int{
//   {0, 0, 0},
//   {0, 0, 0},
//   {0, 1, 0},
// }
// Will flood coord 1,0.
func (s *Surface) fillRows(rows [][]int) {
	for y, row := range rows {
		for x, val := range row {
			if val == 0 {
				continue
			}
			s.Fill(Coord{x, s.height - y - 1})
		}
	}
}
