package utils

import (
	"math"

	geod "github.com/starboard-nz/go-geodesy"
	"github.com/starboard-nz/orb"
	"github.com/starboard-nz/units"
)

// DensifyMultiPolygon inserts points into the multipolygon using the given Model, until the maximum distance between
// planar geometry and the given model is less than the tolerance.
func DensifyMultiPolygon(mp orb.MultiPolygon, tolerance units.Distance, earthModel func(geod.LatLon, ...interface{}) geod.Model) orb.MultiPolygon {
	var dmp orb.MultiPolygon

	for _, polygon := range mp {
		dmp = append(dmp, DensifyPolygon(polygon, tolerance, earthModel))
	}

	return dmp
}

// DensifyPolygon inserts points into the polygon using the given Model, until the maximum distance between
// planar geometry and the given model is less than the tolerance.
func DensifyPolygon(poly orb.Polygon, tolerance units.Distance, earthModel func(geod.LatLon, ...interface{}) geod.Model) orb.Polygon {
	var dp orb.Polygon

	for _, ring := range poly {
		dp = append(dp, DensifyRing(ring, tolerance, earthModel))
	}

	return dp
}

// DensifyRing inserts points into the ring using the given Model, until the maximum distance between
// planar geometry and the given model is less than the tolerance.
func DensifyRing(ring orb.Ring, tolerance units.Distance, earthModel func(geod.LatLon, ...interface{}) geod.Model) orb.Ring {
	if len(ring) < 2 {
		return ring
	}

	lastPoint := ring[len(ring)-1]
	closed := ring[0][0] == lastPoint[0] && ring[0][1] == lastPoint[1]

	points := make([]orb.Point, 0, len(ring))
	dr := orb.Ring(points)
	dr = append(dr, ring[0])
	for i := 1; i < len(ring); i++ {
		ps := DensifySegment(ring[i-1], ring[i], tolerance, earthModel)
		if len(ps) > 1 {
			dr = append(dr, ps[1:]...)
		} else {
			// FIXME log error
		}
	}

	if !closed {
		ps := DensifySegment(lastPoint, ring[0], tolerance, earthModel)
		if len(ps) > 1 {
			dr = append(dr, ps[1:]...)
		} else {
			// FIXME log error
		}
	}

	return dr
}

// DensifySegment inserts intermediate points into the segment p0-p1 using the given Model,
// until the maximum distance between planar geometry and the given model is less than the tolerance.
func DensifySegment(p0, p1 orb.Point, tolerance units.Distance, earthModel func(geod.LatLon, ...interface{}) geod.Model) []orb.Point {
	if SegmentError(p0, p1, earthModel) <= tolerance {
		return []orb.Point{p0, p1}
	}

	ll0 := geod.LatLon{Longitude: geod.Degrees(p0[0]), Latitude: geod.Degrees(p0[1])}
	ll1 := geod.LatLon{Longitude: geod.Degrees(p1[0]), Latitude: geod.Degrees(p1[1])}

	mp := geod.IntermediatePoint(ll0, ll1, 0.5, earthModel)
	omp := orb.Point{float64(mp.Longitude), float64(mp.Latitude)}

	var left, right []orb.Point

	if SegmentError(p0, omp, earthModel) > tolerance {
		left = DensifySegment(p0, omp, tolerance, earthModel)
	} else {
		left = []orb.Point{p0, omp}
	}

	if SegmentError(omp, p1, earthModel) > tolerance {
		right = DensifySegment(omp, p1, tolerance, earthModel)
	} else {
		right = []orb.Point{omp, p1}
	}

	ds := make([]orb.Point, 0, len(left)+len(right)-1)
	ds = append(ds, left...)
	ds = append(ds, right[1:]...)

	return ds
}

// SegmentError calculates the distance between the middle point of a segment calculated using planar geometry
// and using the given Model.
func SegmentError(p0, p1 orb.Point, model geod.EarthModel) units.Distance {
	// following a longitude circle, all supported models are identical
	if p0[0] == p1[0] {
		return 0
	}

	ll0 := geod.LatLon{Latitude: geod.Degrees(p0[1]), Longitude: geod.Degrees(p0[0])}
	ll1 := geod.LatLon{Latitude: geod.Degrees(p1[1]), Longitude: geod.Degrees(p1[0])}

	// halfway point using the given model
	llMid := geod.IntermediatePoint(ll0, ll1, 0.5, model)
	if llMid.Longitude == -180 {
		// distance calculation between -180 and 180 gets confused, so choosing +180 as the antimeridian
		llMid.Longitude = 180
	}

	// planar halfway point (aka Flat Earth™)
	var llFE geod.LatLon

	if math.Abs(p0[0]-p1[0]) < 180 {
		llFE = geod.LatLon{Latitude: geod.Degrees((p0[1] + p1[1]) / 2), Longitude: geod.Degrees((p0[0] + p1[0]) / 2)}
	} else {
		mid := (p0[0]+p1[0])/2 + 180
		if mid > 180 {
			mid -= 360
		} else if mid == -180 {
			mid = 180
		}
		llFE = geod.LatLon{Latitude: geod.Degrees((p0[1] + p1[1]) / 2), Longitude: geod.Degrees(mid)}
	}

	return geod.Distance(llFE, llMid, model)
}
