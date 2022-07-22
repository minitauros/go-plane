package plane

import (
	"testing"
)

func Benchmark_fill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewSurface(11, 11)
		ff := NewFloodFiller(s)
		ff.flood(Coord{0, 0}, Coord{10, 10}, false)
	}
}

func Benchmark_fillAndCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewSurface(11, 11)
		ff := NewFloodFiller(s)
		ff.flood(Coord{0, 0}, Coord{10, 10}, true)
	}
}
