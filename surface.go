package plane

import (
	"sync"
)

// surfaceMap defines which coordinates in a surface are filled.
// The structure is [x][y]value.
type surfaceMap map[int]map[int]int

// Surface represents a surface of a given width and height.
// For width 5 and height 5, the coordinates would range from 0-4x and 0-4y.
// 0,0 is bottom Left.
type Surface struct {
	width  int
	height int
	// surface keeps track of which coordinates are filled.
	surface surfaceMap
	mux     *sync.RWMutex
}

// NewSurface returns a new surface.
func NewSurface(width int, height int) *Surface {
	s := make(surfaceMap, width)
	for i := 0; i < width; i++ {
		s[i] = make(map[int]int, height)
	}
	return &Surface{
		width:   width,
		height:  height,
		surface: s,
		mux:     &sync.RWMutex{},
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
				s.mux.RLock()
				_, ok := s.surface[x][y]
				s.mux.RUnlock()
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
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, coord := range coords {
		delete(s.surface[coord.X], coord.Y)
	}
}

// Fill fills the given coords.
func (s *Surface) Fill(coords ...Coord) {
	for _, coord := range coords {
		if !s.Fits(coord) {
			continue
		}
		// If already filled, don't fill again.
		// We don't want to overwrite the value.
		if s.IsFilled(coord) {
			continue
		}
		s.mux.Lock()
		s.surface[coord.X][coord.Y] = 0
		s.mux.Unlock()
	}
}

// IsFilled returns true if the given coord is filled.
func (s *Surface) IsFilled(coord Coord) bool {
	s.mux.RLock()
	_, ok := s.surface[coord.X][coord.Y]
	s.mux.RUnlock()
	if ok {
		return true
	}
	return coord.X < 0 || coord.Y < 0 || coord.X > s.width-1 || coord.Y > s.height-1
}

// Clone returns a clone of the surface.
func (s *Surface) Clone() *Surface {
	clone := NewSurface(s.width, s.height)
	for coord := range s.EachFilled() {
		clone.forceFill(coord)
	}
	return clone
}

func (s *Surface) fillWithValue(coord Coord, val int) {
	s.Fill(coord)
	s.setValue(coord, val)
}

func (s *Surface) hasValue(coord Coord) bool {
	s.mux.RLock()
	_, ok := s.surface[coord.X][coord.Y]
	s.mux.RUnlock()
	return ok
}

func (s *Surface) getValue(coord Coord) int {
	s.mux.RLock()
	if v, ok := s.surface[coord.X][coord.Y]; ok {
		return v
	}
	s.mux.RUnlock()
	return -1
}

func (s *Surface) setValue(coord Coord, val int) {
	s.mux.Lock()
	s.surface[coord.X][coord.Y] = val
	s.mux.Unlock()
}

// forceFill does a regular fill, but does not check if the coords actually fit on the surface,
// which is slightly faster.
func (s *Surface) forceFill(coords ...Coord) {
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, coord := range coords {
		s.surface[coord.X][coord.Y] = 0
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
