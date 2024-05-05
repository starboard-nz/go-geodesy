package utils

import (
	"errors"
	"fmt"

	geod "github.com/starboard-nz/go-geodesy"
	"github.com/starboard-nz/orb"
	"github.com/starboard-nz/units"
)

var (
	ErrInvalidTolerance = errors.New("invalid value for tolerance - must be positive")
	ErrToleranceTooLow  = errors.New("tolerance too low")
	ErrInternalError    = errors.New("internal error")
	ErrInvalidGeometry  = errors.New("invalid geometry")
)

// NOTE - these densify functions will only work as expected if passing geometries in the normalised (-180 to 180) range.

// DensifyMultiPolygon inserts points into the multipolygon using the given Model, until the maximum distance between
// model and the reference model is less than the tolerance, where model defines the shape of the lines between points
// (e.g. great circle arc or rhumb line).
func DensifyMultiPolygon(mp orb.MultiPolygon, model, refModel geod.EarthModel, tolerance units.Distance) (orb.MultiPolygon, error) {
	var (
		dmp orb.MultiPolygon
		err error
	)

	for _, polygon := range mp {
		dp, err2 := DensifyPolygon(polygon, model, refModel, tolerance)
		if err2 != nil {
			if !errors.Is(err2, ErrToleranceTooLow) {
				return nil, err2
			}

			err = err2
		}

		dmp = append(dmp, dp)
	}

	return dmp, err
}

// DensifyPolygon inserts points into the polygon using the given Model, until the maximum distance between
// planar geometry and the given model is less than the tolerance.
func DensifyPolygon(poly orb.Polygon, model, refModel geod.EarthModel, tolerance units.Distance) (orb.Polygon, error) {
	var (
		dp  orb.Polygon
		err error
	)

	for _, ring := range poly {
		dr, err2 := DensifyRing(ring, model, refModel, tolerance)
		if err2 != nil {
			if !errors.Is(err2, ErrToleranceTooLow) {
				return nil, err2
			}

			err = err2
		}

		dp = append(dp, dr)
	}

	return dp, err
}

// DensifyRing inserts points into the ring using the given Model, until the maximum distance between
// planar geometry and the given model is less than the tolerance.
func DensifyRing(ring orb.Ring, model, refModel geod.EarthModel, tolerance units.Distance) (orb.Ring, error) {
	if len(ring) < 2 {
		return nil, fmt.Errorf("%w: ring has %d points only", ErrInvalidGeometry, len(ring))
	}

	lastPoint := ring[len(ring)-1]
	closed := ring[0][0] == lastPoint[0] && ring[0][1] == lastPoint[1]

	var err error

	points := make([]orb.Point, 0, len(ring))
	dr := orb.Ring(points)
	dr = append(dr, ring[0])
	for i := 1; i < len(ring); i++ {
		ps, err2 := DensifySegment(ring[i-1], ring[i], model, refModel, tolerance)
		if err2 != nil {
			if !errors.Is(err2, ErrToleranceTooLow) {
				return nil, err2
			}

			err = err2
		}

		if len(ps) > 1 {
			dr = append(dr, ps[1:]...)
		} else {
			return nil, ErrInternalError
		}
	}

	if !closed {
		ps, err2 := DensifySegment(lastPoint, ring[0], model, refModel, tolerance)
		if err2 != nil {
			if !errors.Is(err2, ErrToleranceTooLow) {
				return nil, err2
			}

			err = err2
		}

		if len(ps) > 1 {
			dr = append(dr, ps[1:]...)
		} else {
			return nil, ErrInternalError
		}
	}

	return dr, err
}

// DensifySegment inserts intermediate points into the segment p0-p1 using the given Model,
// until the maximum distance between planar geometry and the given model is less than the tolerance.
// If the required tolerance if too low, this function won't exhaust the available memory, but return
// a densified polygon that doesn't meet required tolerance and ErrToleranceTooLow.
func DensifySegment(p0, p1 orb.Point, model, refModel geod.EarthModel, tolerance units.Distance) ([]orb.Point, error) {
	if tolerance.Metre() <= 0 {
		return nil, ErrInvalidTolerance
	}

	ll0 := geod.LatLon{Longitude: geod.Degrees(p0[0]), Latitude: geod.Degrees(p0[1])}
	ll1 := geod.LatLon{Longitude: geod.Degrees(p1[0]), Latitude: geod.Degrees(p1[1])}

	// max 15 deep recursion, allows adding up to 2^14=16364 point per segment, "ought to be enough for anybody"
	return densifySegment(ll0, ll1, p0, p1, 0, 1, model, refModel, tolerance, 15)
}

func densifySegment(ll0, ll1 geod.LatLon, pf, pt orb.Point, from, to float64, model, refModel geod.EarthModel, tolerance units.Distance, recDepth int) ([]orb.Point, error) {
	recDepth -= 1
	mid := (from+to)/2

	mp := geod.IntermediatePoint(ll0, ll1, mid, model)

	llf := geod.LatLon{Latitude: geod.Degrees(pf[1]), Longitude: geod.Degrees(pf[0])}
	llt := geod.LatLon{Latitude: geod.Degrees(pt[1]), Longitude: geod.Degrees(pt[0])}
	refMp := geod.IntermediatePoint(llf, llt, 0.5, refModel)
	e := geod.Distance(mp, refMp, model).Metre()

	if  e.Metre() <= tolerance.Metre() {
		return []orb.Point{pf, pt}, nil
	}

	if recDepth == 0 {
		return []orb.Point{pf, pt}, ErrToleranceTooLow
	}

	var (
		left, right []orb.Point
		err, err2   error
	)

	// middle point (mp) as orb.Point
	omp := orb.Point{float64(mp.Longitude), float64(mp.Latitude)}

	left, err2 = densifySegment(ll0, ll1, pf, omp, from, mid, model, refModel, tolerance, recDepth)
	if err2 != nil {
		if !errors.Is(err2, ErrToleranceTooLow) {
			return nil, err2
		}

		err = err2
	}

	right, err2 = densifySegment(ll0, ll1, omp, pt, mid, to, model, refModel, tolerance, recDepth)
	if err2 != nil {
		if !errors.Is(err2, ErrToleranceTooLow) {
			return nil, err2
		}

		err = err2
	}

	ds := make([]orb.Point, 0, len(left)+len(right)-1)
	ds = append(ds, left...)
	ds = append(ds, right[1:]...)

	return ds, err
}

// SegmentError calculates the distance between the middle point of a segment calculated using planar geometry
// and using the given Model.
func SegmentError(p0, p1 orb.Point, model, refModel geod.EarthModel) units.Distance {
	// following a longitude circle, all supported models are identical
	if p0[0] == p1[0] {
		return units.Metre(0)
	}

	ll0 := geod.LatLon{Latitude: geod.Degrees(p0[1]), Longitude: geod.Degrees(p0[0])}
	ll1 := geod.LatLon{Latitude: geod.Degrees(p1[1]), Longitude: geod.Degrees(p1[0])}

	// halfway point using the given model
	llMid := geod.IntermediatePoint(ll0, ll1, 0.5, model)
	if llMid.Longitude == -180 {
		// distance calculation between -180 and 180 gets confused, so choosing +180 as the antimeridian
		llMid.Longitude = 180
	}

	// reference halfway point
	llRef := geod.IntermediatePoint(ll0, ll1, 0.5, refModel)

	return geod.Distance(llMid, llRef, model)
}
