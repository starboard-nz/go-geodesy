package geod

// Pure Go re-implementation of https://github.com/chrisveness/geodesy

import (
	"github.com/starboard-nz/units"
)

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

/**
 * Library of geodesy functions for operations on a spherical earth model.
 *
 * Includes distances, bearings, destinations, etc, for great circle paths,
 * and other related functions.
 *
 * All calculations are done using simple spherical trigonometric formulae.
 */

import (
	"math"
	"sync"
)

// LatLonSpherical represents a point used for calculations using a spherical Earth model, along great circles
type LatLonSpherical struct {
	ll LatLon
}

// SphericalModel returns a `Model` that wraps geodesy calculations using spherical Earth model along great circles
func SphericalModel(ll LatLon, modelArgs ...interface{}) Model {
	if len(modelArgs) != 0 {
		panic("Invalid number of arguments in call to VincentyModel()")
	}
	return LatLonSpherical{ll: ll}
}

// LatLon converts LatLonSpherical to LatLon
func (lls LatLonSpherical)LatLon() LatLon {
	return lls.ll
}

var earthRadius float64 = 6371000    // metres

// SetEarthRadius can be used to [globally] change the value of Earth's radius (in metres) used
// for spherical Earth calculations (includes rhumb). Default is 6371000m
func SetEarthRadius(r float64) {
	if math.IsNaN(r) {
		panic("Invalid Earth radius specified: NaN")
	}
	if r <= 0 {
		panic("Invalid Earth radius specified, must be positive")
	}
	earthRadius = r
}

// NewLatLonSpherical creates a new LatLonSpherical struct
func NewLatLonSpherical(latitude, longitude float64) LatLonSpherical {
	return LatLonSpherical{
		ll: LatLon{
			Latitude: Wrap90(Degrees(latitude)),
			Longitude: Wrap180(Degrees(longitude)),
		},
	}
}

// ParseLatLonSpherical parses a latitude/longitude point from a variety of formats
// See ParseLatLon for details.
func ParseLatLonSpherical(args ...interface{}) (LatLonSpherical, error) {
	ll, err := ParseLatLon(args)
	if err != nil {
		return LatLonSpherical{}, err
	}
	return LatLonSpherical{ll: ll}, nil
}

// DistanceTo returns the distance along the surface of the earth from `lls` to `dest`.
//
// Uses haversine formula: a = sin²(Δφ/2) + cosφ1·cosφ2 · sin²(Δλ/2); d = 2 · atan2(√a, √(a-1)).
// Use SetEarthRadius() to change the default value.
//
// Argument:
//
// dest  - destination point
//
// Returns the `Distance` between this point and destination point in Distance units.
//
// Examples:
// p1 := geod.NewLatLonSpherical(52.205, 0.119)
// p2 := geod.LatLon{48.857, 2.351}
// d := p1.DistanceTo(p2).Metres()       // 404.3×10³ m
// m := p1.DistanceTo(p2, 3959).Miles()  // 251.2 miles
func (lls LatLonSpherical)DistanceTo(dest LatLon) units.Distance {
        // a = sin²(Δφ/2) + cos(φ1)⋅cos(φ2)⋅sin²(Δλ/2)
        // δ = 2·atan2(√(a), √(1−a))
        // see mathforum.org/library/drmath/view/51879.html for derivation

        R := earthRadius
        φ1 := lls.ll.Latitude.Radians()
	λ1 := lls.ll.Longitude.Radians()
        φ2 := dest.Latitude.Radians()
	λ2 := dest.Longitude.Radians()
        Δφ := φ2 - φ1
        Δλ := λ2 - λ1

        a := math.Sin(Δφ/2) * math.Sin(Δφ/2) + math.Cos(φ1) * math.Cos(φ2) * math.Sin(Δλ/2) * math.Sin(Δλ/2)
        c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
        d := R * c

        return units.Metre(d)
}


// InitialBearingTo returns the initial bearing from `lls` to `dest`.
//
// Argument:
//
// dest  - destination point
//
// Returns the initial bearing in `Degrees` from North (0°..360°)
//
// Example:
// p1 := geod.NewLatLonSpherical(52.205, 0.119)
// p2 := geod.LatLon{48.857, 2.351}
// b1 := p1.InitialBearingTo(p2)    // 156.2°
func (lls LatLonSpherical)InitialBearingTo(dest LatLon) Degrees {
	if lls.ll.Equals(dest) {
		return Degrees(math.NaN())
	}

        // tanθ = sinΔλ⋅cosφ2 / cosφ1⋅sinφ2 − sinφ1⋅cosφ2⋅cosΔλ
        // see mathforum.org/library/drmath/view/55417.html for derivation

        φ1 := lls.ll.Latitude.Radians()
        φ2 := dest.Latitude.Radians()
        Δλ := (dest.Longitude - lls.ll.Longitude).Radians()

        x := math.Cos(φ1) * math.Sin(φ2) - math.Sin(φ1) * math.Cos(φ2) * math.Cos(Δλ)
        y := math.Sin(Δλ) * math.Cos(φ2)
        θ := math.Atan2(y, x)

        bearing := DegreesFromRadians(θ)

        return Wrap360(bearing)
}

// FinalBearingOn returns the final bearing arriving at `dest` from `lls`; the final bearing will
// differ from the initial bearing by varying degrees according to distance and latitude.
//
// Argument:
//
// dest  - destination point
//
// Returns the initial bearing in `Degrees` from North (0°..360°)
//
// Example:
// p1 := geod.NewLatLonSpherical(52.205, 0.119)
// p2 := geod.LatLon{48.857, 2.351}
// b1 := p1.FinalBearingOn(p2)    // 157.9°
func (lls LatLonSpherical)FinalBearingOn(dest LatLon) Degrees {
        // get initial bearing from destination point to this point & reverse it by adding 180°
        bearing := LatLonSpherical{ll: dest}.InitialBearingTo(lls.ll) + 180

        return Wrap360(bearing)
}

// MidPointTo returns the midpoint between `lls` and `dest`
//
// Argument:
//
// dest  - destination point
//
// Returns the middle point
//
// Example:
// p1 := geod.NewLatLonSpherical(52.205, 0.119)
// p2 := geod.LatLon{48.857, 2.351}
// pMid := p1.MidPointTo(p2)    // 50.5363°N, 001.2746°E
func (lls LatLonSpherical)MidPointTo(dest LatLon) LatLon {
        // φm = atan2( sinφ1 + sinφ2, √( (cosφ1 + cosφ2⋅cosΔλ)² + cos²φ2⋅sin²Δλ ) )
        // λm = λ1 + atan2(cosφ2⋅sinΔλ, cosφ1 + cosφ2⋅cosΔλ)
        // midpoint is sum of vectors to two points: mathforum.org/library/drmath/view/51822.html

        φ1 := lls.ll.Latitude.Radians()
        λ1 := lls.ll.Longitude.Radians()
        φ2 := dest.Latitude.Radians()
        Δλ := (dest.Longitude - lls.ll.Longitude).Radians()

        // get cartesian coordinates for the two points
        A := Cartesian{ X: math.Cos(φ1), Y: 0, Z: math.Sin(φ1) }    // place point A on prime meridian y=0
        B := Cartesian{ X: math.Cos(φ2) * math.Cos(Δλ), Y: math.Cos(φ2) * math.Sin(Δλ), Z: math.Sin(φ2) }

        // vector to midpoint is sum of vectors to two points (no need to normalise)
        C := Cartesian{ X: A.X + B.X, Y: A.Y + B.Y, Z: A.Z + B.Z }

        φm := math.Atan2(C.Z, math.Sqrt(C.X * C.X + C.Y * C.Y))
        λm := λ1 + math.Atan2(C.Y, C.X)

        lat := DegreesFromRadians(φm)
        lon := DegreesFromRadians(λm)

        return LatLon{Latitude: Wrap90(lat), Longitude: Wrap180(lon)}
}

// IntermediatePointTo returns the point at the given fraction between `lls` and `dest`.
//
// Arguments:
//
// dest  - destination point
// fraction - Fraction between the two points (0 = `lls`, 1 = `dest`)
//
// Returns the intermediate point.
//
// Example:
// p1 := geod.NewLatLonSpherical(52.205, 0.119)
// p2 := geod.LatLon{48.857, 2.351}
// pInt := p1.IntermediatePointTo(p2, 0.25)    // 51.3721°N, 000.7073°E
func (lls LatLonSpherical)IntermediatePointTo(dest LatLon, fraction float64) LatLon {
	if lls.ll.Equals(dest) {
		return lls.ll
	}

        φ1 := lls.ll.Latitude.Radians()
	λ1 := lls.ll.Longitude.Radians()
        φ2 := dest.Latitude.Radians()
	λ2 := dest.Longitude.Radians()

        // distance between points
        Δφ := φ2 - φ1
        Δλ := λ2 - λ1
        a := math.Sin(Δφ/2) * math.Sin(Δφ/2) + math.Cos(φ1) * math.Cos(φ2) * math.Sin(Δλ/2) * math.Sin(Δλ/2)
        δ := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

        A := math.Sin((1 - fraction)*δ) / math.Sin(δ)
        B := math.Sin(fraction * δ) / math.Sin(δ)

        x := A * math.Cos(φ1) * math.Cos(λ1) + B * math.Cos(φ2) * math.Cos(λ2)
        y := A * math.Cos(φ1) * math.Sin(λ1) + B * math.Cos(φ2) * math.Sin(λ2)
        z := A * math.Sin(φ1) + B * math.Sin(φ2)

        φ3 := math.Atan2(z, math.Sqrt(x*x + y*y))
        λ3 := math.Atan2(y, x)

        lat := DegreesFromRadians(φ3)
        lon := DegreesFromRadians(λ3)

        return LatLon{Latitude: Wrap90(lat), Longitude: Wrap180(lon)}
}

// IntermediatePointsTo returns the points at the given fractions between `lls` and `dest`.
//
// Arguments:
//
// dest  - destination point
// fraction - Slice of fractions between the two points (0 = `lls`, 1 = `dest`)
//
// Returns an intermediate point for each fraction
//
// Example:
// p1 := geod.NewLatLonSpherical(52.205, 0.119)
// p2 := geod.LatLon{48.857, 2.351}
// pInt := p1.IntermediatePointsTo(p2, []float64{0.25, 0.5, 0.75})
func (lls LatLonSpherical)IntermediatePointsTo(dest LatLon, fractions []float64) []LatLon {
	waitGroup := &sync.WaitGroup{}

	points := make([]LatLon, len(fractions))
	for i, fraction := range(fractions) {
		waitGroup.Add(1)
		go func(i int, fraction float64) {
			points[i] = lls.IntermediatePointTo(dest, fraction)
			waitGroup.Done()
		} (i, fraction)
	}

	// wait for all goroutines to finish
	waitGroup.Wait()

	return points
}


// DestinationPoint returns the destination point from `lls` having travelled the given distance on the
// given initial bearing (bearing normally varies around path followed).
//
// Arguments:
//
// distance - Distance travelled in metres
// bearing - Initial bearing in `Degrees` from North
//
// Returns the destination point.
//
// Example:
// p1 := geod.NewLatLonSpherical(51.47788, -0.00147)
// p2 := p1.DestinationPoint(7794, geod.Degrees(300.7)) // 51.5136°N, 000.0983°W
func (lls LatLonSpherical)DestinationPoint(distance float64, bearing Degrees) LatLon {
        // sinφ2 = sinφ1⋅cosδ + cosφ1⋅sinδ⋅cosθ
        // tanΔλ = sinθ⋅sinδ⋅cosφ1 / cosδ−sinφ1⋅sinφ2
        // see mathforum.org/library/drmath/view/52049.html for derivation

        δ := distance / earthRadius     // angular distance in radians
        θ := bearing.Radians()

        φ1 := lls.ll.Latitude.Radians()
	λ1 := lls.ll.Longitude.Radians()

        sinφ2 := math.Sin(φ1) * math.Cos(δ) + math.Cos(φ1) * math.Sin(δ) * math.Cos(θ)
        φ2 := math.Asin(sinφ2)
        y := math.Sin(θ) * math.Sin(δ) * math.Cos(φ1)
        x := math.Cos(δ) - math.Sin(φ1) * sinφ2
        λ2 := λ1 + math.Atan2(y, x)

        lat := DegreesFromRadians(φ2)
        lon := DegreesFromRadians(λ2)

        return LatLon{Latitude: Wrap90(lat), Longitude: Wrap180(lon)}
}

// Intersection returns the point of intersection of two paths defined by point and bearing.
//
// Arguments:
//
// bearing1 - Initial bearing in `Degrees` from North from `lls`
// lls2 - Second point
// bearing2 - Initial bearing in `Degrees` from North from `lls2`
//
// Returns the point of intersection of the 2 paths.
// If the intersection point cannot be calculated (e.g. infinite intersections) the returned point
// has NaN as Latitude and Longitude.
//
// Example:
// p1 := geod.NewLatLonSpherical(51.8853, 0.2545)
// brng1 := geod.Degrees(108.547)
// p2 := geod.LatLon{49.0034, 2.5735}
// brng2 := geod.Degrees(32.435)
// pInt := p1.Intersection(brng1, p2, brng2) // 50.9078°N, 004.5084°E
func (lls LatLonSpherical)Intersection(bearing1 Degrees, ll2 LatLon, bearing2 Degrees) LatLon {
	const π = math.Pi
	ε := math.Nextafter(1, 2) - 1
	
        // see www.edwilliams.org/avform.htm#Intersection

        φ1 := lls.ll.Latitude.Radians()
	λ1 := lls.ll.Longitude.Radians()
        φ2 := ll2.Latitude.Radians()
	λ2 := ll2.Longitude.Radians()
        θ13 := bearing1.Radians()
	θ23 := bearing2.Radians()
        Δφ := φ2 - φ1
	Δλ := λ2 - λ1

        // angular distance p1-p2
        δ12 := 2 * math.Asin(math.Sqrt(math.Sin(Δφ/2) * math.Sin(Δφ/2) +
		math.Cos(φ1) * math.Cos(φ2) * math.Sin(Δλ/2) * math.Sin(Δλ/2)))
        if math.Abs(δ12) < ε {
		return lls.ll  // coincident points
	}

        // initial/final bearings between points
        cosθa := (math.Sin(φ2) - math.Sin(φ1)*math.Cos(δ12)) / (math.Sin(δ12)*math.Cos(φ1))
        cosθb := (math.Sin(φ1) - math.Sin(φ2)*math.Cos(δ12)) / (math.Sin(δ12)*math.Cos(φ2))
        θa := math.Acos(math.Min(math.Max(cosθa, -1), 1))       // protect against rounding errors
        θb := math.Acos(math.Min(math.Max(cosθb, -1), 1))       // protect against rounding errors

        θ12 := θa
	if math.Sin(λ2-λ1) <= 0 {
		θ12 = 2*π-θa
	}
        θ21 := θb
	if math.Sin(λ2-λ1)>0 {
		θ21 = 2*π-θb
	}

        α1 := θ13 - θ12    // angle 2-1-3
        α2 := θ21 - θ23    // angle 1-2-3

        if math.Sin(α1) == 0 && math.Sin(α2) == 0 {
		return LatLon{Latitude: Degrees(math.NaN()), Longitude: Degrees(math.NaN())}  // infinite intersections
	}
        if math.Sin(α1) * math.Sin(α2) < 0 {
		return LatLon{Latitude: Degrees(math.NaN()), Longitude: Degrees(math.NaN())}  // ambiguous intersection (antipodal?)
	}

        cosα3 := -math.Cos(α1)*math.Cos(α2) + math.Sin(α1)*math.Sin(α2)*math.Cos(δ12)

        δ13 := math.Atan2(math.Sin(δ12)*math.Sin(α1)*math.Sin(α2), math.Cos(α2) + math.Cos(α1)*cosα3)

        φ3 := math.Asin(math.Min(math.Max(math.Sin(φ1)*math.Cos(δ13) + math.Cos(φ1)*math.Sin(δ13)*math.Cos(θ13), -1), 1))

        Δλ13 := math.Atan2(math.Sin(θ13)*math.Sin(δ13)*math.Cos(φ1), math.Cos(δ13) - math.Sin(φ1)*math.Sin(φ3))
        λ3 := λ1 + Δλ13

        lat := DegreesFromRadians(φ3)
        lon := DegreesFromRadians(λ3)

        return LatLon{Latitude: Wrap90(lat), Longitude: Wrap180(lon)}
}

