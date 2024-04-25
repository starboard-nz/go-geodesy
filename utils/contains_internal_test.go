package utils

import (
	"math/rand"
	"testing"

	geod "github.com/starboard-nz/go-geodesy"
	"github.com/starboard-nz/orb"
)

func BenchmarkRayIntersect(b *testing.B) {
	const N = 100000
	s := make([]orb.Point, N)
	e := make([]orb.Point, N)
	p := make([]orb.Point, N)
	for i := 0; i < N; i++ {
		s[i] = orb.Point{
			rand.Float64() * 100, // nolint:gosec
			rand.Float64() * 100, // nolint:gosec
		}
		e[i] = orb.Point{
			rand.Float64() * 100, // nolint:gosec
			rand.Float64() * 100, // nolint:gosec
		}
		p[i] = orb.Point{
			rand.Float64() * 100, // nolint:gosec
			rand.Float64() * 100, // nolint:gosec
		}
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		rayIntersect(p[n%N], s[n%N], e[n%N], geod.RhumbModel)
	}
}
