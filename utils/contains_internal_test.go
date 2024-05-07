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
	start := geod.LatLon{Latitude: 0, Longitude: 0}
	end := geod.LatLon{Latitude: -50, Longitude: 60}
	point := geod.LatLon{Latitude: -10, Longitude: 55}
	bse := geod.InitialBearing(start, end, geod.SphericalModel)
	bsp := geod.InitialBearing(start, point, geod.SphericalModel)
	t.Logf("s -> e initial bearing: %v", bse)
	t.Logf("s -> p initial bearing: %v", bsp)

	// Edge case tests

	s := orb.Point{0, 0}
	e := orb.Point{5, 5}
	p := orb.Point{5, 2}  // Shares same x as end point
	p2 := orb.Point{0, 5} // Shares same x as start point
	p3 := orb.Point{5, 5} // Identical to end point
	p4 := orb.Point{0, 0} // Identical to start point

	i, on := rayIntersect(p, s, e, geod.RhumbModel)
	if i { // Test shared x values don't trigger contains for both start and end of segment
		t.Errorf("point sharing same x as end and being under the test segment should not intersect")
	}
	if on {
		t.Errorf("on should be false")
	}

	i, on = rayIntersect(p2, s, e, geod.RhumbModel)
	if i {
		t.Errorf("point sharing same x as start point and being below the segment should intersect")
	}
	if on {
		t.Errorf("on should be false")
	}

	i, on = rayIntersect(p3, s, e, geod.RhumbModel)
	if i {
		t.Errorf("intersection should be false as test point is identical to end point")
	}
	if !on {
		t.Errorf("on should be true as test point is identical to end point")
	}

	i, on = rayIntersect(p4, s, e, geod.RhumbModel)
	if i {
		t.Errorf("intersection should be false as test point is identical to end point")
	}
	if !on {
		t.Errorf("on should be true as test point is identical to start point")
	}

	s = orb.Point{0, 5}
	e = orb.Point{5, 0}
	p = orb.Point{5, -2} // Shares same x as end point
	p2 = orb.Point{0, 0} // Shares same x as start point
	p3 = orb.Point{5, 0} // Identical to end point
	p4 = orb.Point{0, 5} // Identical to start point

	i, on = rayIntersect(p, s, e, geod.RhumbModel)
	if i {
		t.Errorf("point sharing same x as end and being under the test segment should not intersect")
	}
	if on {
		t.Errorf("on should be false")
	}

	i, on = rayIntersect(p2, s, e, geod.RhumbModel)
	if !i {
		t.Errorf("point sharing same x as start point and being below the segment should intersect")
	}
	if on {
		t.Errorf("on should be false")
	}

	i, on = rayIntersect(p3, s, e, geod.RhumbModel)
	if i {
		t.Errorf("intersection should be false as test point is identical to end point")
	}
	if !on {
		t.Errorf("on should be true as test point is identical to end point")
	}

	i, on = rayIntersect(p4, s, e, geod.RhumbModel)
	if i {
		t.Errorf("intersection should be false as test point is identical to end point")
	}
	if !on {
		t.Errorf("on should be true as test point is identical to start point")
	}
}
