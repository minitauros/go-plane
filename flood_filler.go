package plane

import (
	"sync"
)

const maxConcurrency = 8

// FloodFiller is a flood filler.
type FloodFiller struct {
	s      *Surface
	concCh chan struct{}
}

// NewFloodFiller returns a new flood filler.
func NewFloodFiller(surface *Surface) *FloodFiller {
	return &FloodFiller{
		s:      surface,
		concCh: make(chan struct{}, maxConcurrency),
	}
}

// Fill starts a flood fill from `base`, starting the flood at `startAt`.
// It returns the number of coords that were filled.
// It does not fill `base`.
func (f *FloodFiller) Fill(base, startAt Coord) int {
	if !base.ConnectsTo(startAt) {
		return 0
	}
	numFilledAtStart := f.s.CountFilled()
	f.fill(base, startAt)
	numFilledAtEnd := f.s.CountFilled()
	return numFilledAtEnd - numFilledAtStart
}

// CanReach returns true if a path can be made through unfilled coords from `base` to `target`.
func (f *FloodFiller) CanReach(base, target Coord) bool {
	return f.canReach(base, target)
}

// CanReachWhenStartingFillAt returns true if a path can be made through unfilled coords from `base` to `target`,
// given that the flood fill must start at `startAt`.
func (f *FloodFiller) CanReachWhenStartingFillAt(base, target, startAt Coord) bool {
	return f.canReach(base, target, startAt)
}

// CountSteps returns the smallest number of steps that can be taken to reach `target` from `base`.
func (f *FloodFiller) CountSteps(base, target Coord) int {
	if !f.canReach(base, target) {
		return -1
	}
	return f.s.getValue(target)
}

func (f *FloodFiller) canReach(base, target Coord, allowedStarts ...Coord) bool {
	filledAroundBefore := f.s.getCoordsFilledAround(target)
	for _, d := range base.GetDirectionsTo(target) {
		coordInDirection := base.GetCoordInDirection(d)
		if len(allowedStarts) > 0 {
			var mayStartInThisDirection bool
			for _, allowedStart := range allowedStarts {
				if coordInDirection == allowedStart {
					mayStartInThisDirection = true
					break
				}
			}
			if !mayStartInThisDirection {
				continue
			}
		}
		f.fill(base, coordInDirection)
	}
	filledAroundAfter := f.s.getCoordsFilledAround(target)
	return len(filledAroundAfter) > len(filledAroundBefore)
}

func (f *FloodFiller) fill(base, start Coord) {
	// fill base coord so that the flood cannot escape.
	f.s.Fill(base)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// We start with number of steps taken = 1 because the starting coord is already 1 step away from the base coord, and we are also exploring that one.
	go f.explore(start, start.GetDirectionsTo(base)[0], wg, true, 1)
	wg.Wait()

	f.s.Remove(base)
}

func (f *FloodFiller) explore(
	target Coord,
	skipDirection Direction,
	wg *sync.WaitGroup,
	concurrent bool,
	numStepsTaken int,
) {
	numStepsTaken++
	defer wg.Done()
	if concurrent {
		f.concCh <- struct{}{}
		defer func() {
			<-f.concCh
		}()
	}
	if f.s.IsFilled(target) {
		if !f.s.hasValue(target) {
			f.s.setValue(target, numStepsTaken)
		}
		return
	} else {
		f.s.fillWithValue(target, numStepsTaken)
	}
	for _, d := range GetAllDirections() {
		if d == skipDirection {
			continue
		}
		next := target.GetCoordInDirection(d)
		if !f.s.IsFilled(next) {
			wg.Add(1)
			goingInSameDirection := d == skipDirection.opposite()
			if goingInSameDirection {
				f.explore(next, d.opposite(), wg, false, numStepsTaken)
			} else {
				go f.explore(next, d.opposite(), wg, true, numStepsTaken)
			}
		}
	}
}
