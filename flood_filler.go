package plane

// FloodFiller is a flood filler.
type FloodFiller struct {
	s *Surface
}

// NewFloodFiller returns a new flood filler.
func NewFloodFiller(surface *Surface) *FloodFiller {
	return &FloodFiller{
		s: surface,
	}
}

// Flood starts a flood fill from `base`, starting the flood at `startAt`.
// It returns the number of coords that were filled.
// It does not flood `base`.
func (f *FloodFiller) Flood(base, startAt Coord) Coords {
	if !base.ConnectsTo(startAt) {
		return Coords{}
	}
	coordsFilled := f.flood(base, startAt, false)
	return coordsFilled
}

// CanReach returns true if a path can be made through unfilled coords from `base` to `target`.
func (f *FloodFiller) CanReach(base, target Coord) bool {
	return f.canReach(base, target, false)
}

// CanReachWhenStartingFillAt returns true if a path can be made through unfilled coords from `base` to `target`,
// given that the flood fill must start at `startAt`.
func (f *FloodFiller) CanReachWhenStartingFillAt(base, target, startAt Coord) bool {
	return f.canReach(base, target, false, startAt)
}

// CountSteps returns the smallest number of steps that can be taken to reach `target` from `base`.
func (f *FloodFiller) CountSteps(base, target Coord) int {
	if !f.canReach(base, target, true) {
		return -1
	}
	return f.s.getDistance(target)
}

func (f *FloodFiller) canReach(base, target Coord, countSteps bool, allowedStarts ...Coord) bool {
	distanceBefore := f.s.getDistance(target)
	filledAroundBefore := f.s.getCoordsFilledAround(target)
	for _, d := range GetAllDirections() {
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
		f.flood(base, coordInDirection, countSteps)
	}
	if countSteps {
		return distanceBefore != f.s.getDistance(target)
	}
	filledAroundAfter := f.s.getCoordsFilledAround(target)
	return len(filledAroundAfter) > len(filledAroundBefore)
}

func (f *FloodFiller) flood(base, start Coord, countSteps bool) Coords {
	// Flood base coord so that the flood cannot escape.
	f.s.Fill(base)

	filled := &Coords{}
	if countSteps {
		f.exploreDistance(start, start.GetDirectionsTo(base)[0], 0)
	} else {
		f.explore(start, start.GetDirectionsTo(base)[0], filled)
	}

	f.s.Remove(base)
	return *filled
}

func (f *FloodFiller) explore(
	target Coord,
	comingFromDirection Direction,
	filled *Coords,
) {
	if f.s.IsFilled(target) {
		return
	}
	f.s.Fill(target)
	*filled = append(*filled, target)
	for _, d := range GetAllDirections() {
		if d == comingFromDirection {
			continue
		}
		f.explore(target.GetCoordInDirection(d), d.Opposite(), filled)
	}
}

// exploreDistance is like explore but also counts the number of steps it takes to reach positions.
// It is a bit more inefficient because it will re-explore coords that were already reached,
// but that can be explored by a different call/routine in fewer steps.
func (f *FloodFiller) exploreDistance(
	target Coord,
	comingFromDirection Direction,
	numStepsTaken int,
) {
	numStepsTaken++
	if !f.s.Fits(target) {
		return
	}
	curVal, exists := f.s.getValue(target)
	if !exists || curVal.distance == 0 || curVal.distance > numStepsTaken {
		f.s.setDistance(target, numStepsTaken)
	} else if exists && curVal.distance <= numStepsTaken {
		return
	}
	if f.s.IsFilled(target) {
		return
	}
	for _, d := range GetAllDirections() {
		if d == comingFromDirection {
			continue
		}
		f.exploreDistance(target.GetCoordInDirection(d), d.Opposite(), numStepsTaken)
	}
}
