package plane

// Direction exists so that the human brain can more easily interpret what is happening.
// It is not the most performant way of doing this, but then again, when this package was made, that was not a priority.
// Consider removing it in the future.
type Direction string

const (
	Top   Direction = "top"
	Right Direction = "right"
	Bot   Direction = "bot"
	Left  Direction = "left"
)

// IsVertical returns true if the direction is vertical.
func (d Direction) IsVertical() bool {
	return d == Top || d == Bot
}

// IsHorizontal returns true if the direction is horizontal.
func (d Direction) IsHorizontal() bool {
	return d == Left || d == Right
}

// Opposite returns the opposite direction.
func (d Direction) Opposite() Direction {
	switch d {
	case Top:
		return Bot
	case Bot:
		return Top
	case Right:
		return Left
	case Left:
		return Right
	}
	return Bot
}

// NextClockwise returns the next clockwise direction.
func (d Direction) NextClockwise() Direction {
	switch d {
	case Top:
		return Right
	case Bot:
		return Left
	case Right:
		return Bot
	case Left:
		return Top
	}
	return Top
}

// NextCounterClockwise returns the next counter clockwise direction.
func (d Direction) NextCounterClockwise() Direction {
	switch d {
	case Top:
		return Left
	case Bot:
		return Right
	case Right:
		return Top
	case Left:
		return Bot
	}
	return Top
}

var allDirections = []Direction{Top, Right, Bot, Left}

// GetAllDirections returns all available directions.
func GetAllDirections() []Direction {
	return allDirections
}
