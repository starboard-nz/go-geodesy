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


func TestRayInterSect(t *testing.T) {
	s := geod.LatLon{Latitude: 0, Longitude: 0}
	e := geod.LatLon{Latitude: -50, Longitude: 60}
	p := geod.LatLon{Latitude: -10, Longitude: 55}
	bse := geod.InitialBearing(s, e, geod.SphericalModel)
	bsp := geod.InitialBearing(s, p, geod.SphericalModel)
	t.Logf("s -> e initial bearing: %v", bse)
	t.Logf("s -> p initial bearing: %v", bsp)
}
