package utils

import (
	"math"

	geod "github.com/starboard-nz/go-geodesy"
	"github.com/starboard-nz/orb"
)

// RingContains returns true if the point is inside the ring.
// Points on the boundary are considered in.
func RingContains(r orb.Ring, point orb.Point, model geod.EarthModel) bool {
	if !r.Bound().Contains(point) {
		return false
	}

	c, on := rayIntersect(point, r[0], r[len(r)-1], model)
	if on {
		return true
	}

	for i := 0; i < len(r)-1; i++ {
		inter, on := rayIntersect(point, r[i], r[i+1], model)
		if on {
			return true
		}

		if inter {
			c = !c
		}
	}

	return c
}

// PolygonContains checks if the point is within the polygon.
// Points on the boundary are considered in.
func PolygonContains(p orb.Polygon, point orb.Point, model geod.EarthModel) bool {
	if !RingContains(p[0], point, model) {
		return false
	}

	for i := 1; i < len(p); i++ {
		if RingContains(p[i], point, model) {
			return false
		}
	}

	return true
}

// MultiPolygonContains checks if the point is within the multi-polygon.
// Points on the boundary are considered in.
func MultiPolygonContains(mp orb.MultiPolygon, point orb.Point, model geod.EarthModel) bool {
	for _, p := range mp {
		if PolygonContains(p, point, model) {
			return true
		}
	}

	return false
}

// RingWithBoundContains returns true if the point is inside a ring with the given bound.
// This is an optimization of RingContains that avoids re-calculating the bound for each
// point that is tested.
// Points on the boundary are considered in.
func RingWithBoundContains(r orb.Ring, bound orb.Bound, point orb.Point, model geod.EarthModel) bool {

	if bound.IsZero() || bound.IsEmpty() {
		bound = r.Bound()
	}

	if !bound.Contains(point) {
		return false
	}

	c, on := rayIntersect(point, r[0], r[len(r)-1], model)
	if on {
		return true
	}

	for i := 0; i < len(r)-1; i++ {
		inter, on := rayIntersect(point, r[i], r[i+1], model)
		if on {
			return true
		}

		if inter {
			c = !c
		}
	}

	return c
}

// PolygonWithBoundContains checks if the point is within the polygon with the given bounds.
// The bounds can be calculated using PolygonBoundsFromPolygon().
// This is an optimization of PolygonContains that avoids re-calculating the bounds for each point
// that is tested.
// Points on the boundary are considered in.
func PolygonWithBoundContains(poly orb.Polygon, bounds orb.PolygonBounds, point orb.Point, model geod.EarthModel) bool {
	if bounds == nil {
		bounds = orb.PolygonBoundsFromPolygon(poly)
	}

	if !RingWithBoundContains(poly[0], bounds[0], point, model) {
		return false
	}

	for i := 1; i < len(poly); i++ {
		if RingWithBoundContains(poly[i], bounds[i], point, model) {
			return false
		}
	}

	return true
}

// MultiPolygonWithBoundContains checks if the point is within the multi-polygon with the given bounds.
// The multiBounds can be calculated using MultiPolygonBoundsFromMultiPolygon().
// This is an optimization of MultiPolygonContains that avoids re-calculating the bounds for each point
// that is tested.
// Points on the boundary are considered in.
func MultiPolygonWithBoundContains(mp orb.MultiPolygon, multiBounds orb.MultiPolygonBounds, point orb.Point, model geod.EarthModel) bool {
	if multiBounds == nil {
		multiBounds = orb.MultiPolygonBoundsFromMultiPolygon(mp)
	}

	for i, poly := range mp {
		if PolygonWithBoundContains(poly, multiBounds[i], point, model) {
			return true
		}
	}

	return false
}

// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func rayIntersect(p, s, e orb.Point, model geod.EarthModel) (intersects, on bool) {
	if s[0] > e[0] {
		s, e = e, s
	}

	if p[0] == s[0] {
		if p[1] == s[1] {
			// p == start
			return false, true
		} else if s[0] == e[0] {
			// vertical segment (s -> e)
			// return true if within the line, check to see if start or end is greater.
			if s[1] > e[1] && s[1] >= p[1] && p[1] >= e[1] {
				return false, true
			}

			if e[1] > s[1] && e[1] >= p[1] && p[1] >= s[1] {
				return false, true
			}
		}

		// Move the y coordinate to deal with degenerate case
		p[0] = math.Nextafter(p[0], math.Inf(1))
	} else if p[0] == e[0] {
		if p[1] == e[1] {
			// matching the end point
			return false, true
		}

		p[0] = math.Nextafter(p[0], math.Inf(1))
	}

	if p[0] < s[0] || p[0] > e[0] {
		return false, false
	}

	if s[1] > e[1] {
		if p[1] > s[1] {
			return false, false
		} else if p[1] < e[1] {
			return true, false
		}
	} else {
		if p[1] > e[1] {
			return false, false
		} else if p[1] < s[1] {
			return true, false
		}
	}

	bs := geod.InitialBearing(
		geod.LatLon{Latitude: geod.Degrees(s[1]), Longitude: geod.Degrees(s[0])},
		geod.LatLon{Latitude: geod.Degrees(e[1]), Longitude: geod.Degrees(e[0])},
		model)
	bp := geod.InitialBearing(
		geod.LatLon{Latitude: geod.Degrees(s[1]), Longitude: geod.Degrees(s[0])},
		geod.LatLon{Latitude: geod.Degrees(p[1]), Longitude: geod.Degrees(p[0])},
		model)

	if bs == bp {
		return false, true
	}

	return bs <= bp, false
}
