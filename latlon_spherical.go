package geod

// Pure Go re-implementation of https://github.com/chrisveness/geodesy

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
// Uses haversine formula: a = sin??(????/2) + cos??1??cos??2 ?? sin??(????/2); d = 2 ?? atan2(???a, ???(a-1)).
// Use SetEarthRadius() to change the default value.
//
// Argument:
//
// dest  - destination point
//
// Returns the `Distance` between this point and destination point in DistanceUnits
//
// Examples:
// p1 := geod.NewLatLonSpherical(52.205, 0.119)
// p2 := geod.LatLon{48.857, 2.351}
// d := p1.DistanceTo(p2).Metres()       // 404.3??10?? m
// m := p1.DistanceTo(p2, 3959).Miles()  // 251.2 miles
func (lls LatLonSpherical)DistanceTo(dest LatLon) DistanceUnits {
        // a = sin??(????/2) + cos(??1)???cos(??2)???sin??(????/2)
        // ?? = 2??atan2(???(a), ???(1???a))
        // see mathforum.org/library/drmath/view/51879.html for derivation

        R := earthRadius
        ??1 := lls.ll.Latitude.Radians()
	??1 := lls.ll.Longitude.Radians()
        ??2 := dest.Latitude.Radians()
	??2 := dest.Longitude.Radians()
        ???? := ??2 - ??1
        ???? := ??2 - ??1

        a := math.Sin(????/2) * math.Sin(????/2) + math.Cos(??1) * math.Cos(??2) * math.Sin(????/2) * math.Sin(????/2)
        c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
        d := R * c

        return DistanceUnits(d)
}


// InitialBearingTo returns the initial bearing from `lls` to `dest`.
//
// Argument:
//
// dest  - destination point
//
// Returns the initial bearing in `Degrees` from North (0??..360??)
//
// Example:
// p1 := geod.NewLatLonSpherical(52.205, 0.119)
// p2 := geod.LatLon{48.857, 2.351}
// b1 := p1.InitialBearingTo(p2)    // 156.2??
func (lls LatLonSpherical)InitialBearingTo(dest LatLon) Degrees {
	if lls.ll.Equals(dest) {
		return Degrees(math.NaN())
	}

        // tan?? = sin???????cos??2 / cos??1???sin??2 ??? sin??1???cos??2???cos????
        // see mathforum.org/library/drmath/view/55417.html for derivation

        ??1 := lls.ll.Latitude.Radians()
        ??2 := dest.Latitude.Radians()
        ???? := (dest.Longitude - lls.ll.Longitude).Radians()

        x := math.Cos(??1) * math.Sin(??2) - math.Sin(??1) * math.Cos(??2) * math.Cos(????)
        y := math.Sin(????) * math.Cos(??2)
        ?? := math.Atan2(y, x)

        bearing := DegreesFromRadians(??)

        return Wrap360(bearing)
}

// FinalBearingOn returns the final bearing arriving at `dest` from `lls`; the final bearing will
// differ from the initial bearing by varying degrees according to distance and latitude.
//
// Argument:
//
// dest  - destination point
//
// Returns the initial bearing in `Degrees` from North (0??..360??)
//
// Example:
// p1 := geod.NewLatLonSpherical(52.205, 0.119)
// p2 := geod.LatLon{48.857, 2.351}
// b1 := p1.FinalBearingOn(p2)    // 157.9??
func (lls LatLonSpherical)FinalBearingOn(dest LatLon) Degrees {
        // get initial bearing from destination point to this point & reverse it by adding 180??
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
// pMid := p1.MidPointTo(p2)    // 50.5363??N, 001.2746??E
func (lls LatLonSpherical)MidPointTo(dest LatLon) LatLon {
        // ??m = atan2( sin??1 + sin??2, ???( (cos??1 + cos??2???cos????)?? + cos????2???sin?????? ) )
        // ??m = ??1 + atan2(cos??2???sin????, cos??1 + cos??2???cos????)
        // midpoint is sum of vectors to two points: mathforum.org/library/drmath/view/51822.html

        ??1 := lls.ll.Latitude.Radians()
        ??1 := lls.ll.Longitude.Radians()
        ??2 := dest.Latitude.Radians()
        ???? := (dest.Longitude - lls.ll.Longitude).Radians()

        // get cartesian coordinates for the two points
        A := Cartesian{ X: math.Cos(??1), Y: 0, Z: math.Sin(??1) }    // place point A on prime meridian y=0
        B := Cartesian{ X: math.Cos(??2) * math.Cos(????), Y: math.Cos(??2) * math.Sin(????), Z: math.Sin(??2) }

        // vector to midpoint is sum of vectors to two points (no need to normalise)
        C := Cartesian{ X: A.X + B.X, Y: A.Y + B.Y, Z: A.Z + B.Z }

        ??m := math.Atan2(C.Z, math.Sqrt(C.X * C.X + C.Y * C.Y))
        ??m := ??1 + math.Atan2(C.Y, C.X)

        lat := DegreesFromRadians(??m)
        lon := DegreesFromRadians(??m)

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
// pInt := p1.IntermediatePointTo(p2, 0.25)    // 51.3721??N, 000.7073??E
func (lls LatLonSpherical)IntermediatePointTo(dest LatLon, fraction float64) LatLon {
	if lls.ll.Equals(dest) {
		return lls.ll
	}

        ??1 := lls.ll.Latitude.Radians()
	??1 := lls.ll.Longitude.Radians()
        ??2 := dest.Latitude.Radians()
	??2 := dest.Longitude.Radians()

        // distance between points
        ???? := ??2 - ??1
        ???? := ??2 - ??1
        a := math.Sin(????/2) * math.Sin(????/2) + math.Cos(??1) * math.Cos(??2) * math.Sin(????/2) * math.Sin(????/2)
        ?? := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

        A := math.Sin((1 - fraction)*??) / math.Sin(??)
        B := math.Sin(fraction * ??) / math.Sin(??)

        x := A * math.Cos(??1) * math.Cos(??1) + B * math.Cos(??2) * math.Cos(??2)
        y := A * math.Cos(??1) * math.Sin(??1) + B * math.Cos(??2) * math.Sin(??2)
        z := A * math.Sin(??1) + B * math.Sin(??2)

        ??3 := math.Atan2(z, math.Sqrt(x*x + y*y))
        ??3 := math.Atan2(y, x)

        lat := DegreesFromRadians(??3)
        lon := DegreesFromRadians(??3)

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
// p2 := p1.DestinationPoint(7794, geod.Degrees(300.7)) // 51.5136??N, 000.0983??W
func (lls LatLonSpherical)DestinationPoint(distance float64, bearing Degrees) LatLon {
        // sin??2 = sin??1???cos?? + cos??1???sin?????cos??
        // tan???? = sin?????sin?????cos??1 / cos?????sin??1???sin??2
        // see mathforum.org/library/drmath/view/52049.html for derivation

        ?? := distance / earthRadius     // angular distance in radians
        ?? := bearing.Radians()

        ??1 := lls.ll.Latitude.Radians()
	??1 := lls.ll.Longitude.Radians()

        sin??2 := math.Sin(??1) * math.Cos(??) + math.Cos(??1) * math.Sin(??) * math.Cos(??)
        ??2 := math.Asin(sin??2)
        y := math.Sin(??) * math.Sin(??) * math.Cos(??1)
        x := math.Cos(??) - math.Sin(??1) * sin??2
        ??2 := ??1 + math.Atan2(y, x)

        lat := DegreesFromRadians(??2)
        lon := DegreesFromRadians(??2)

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
// pInt := p1.Intersection(brng1, p2, brng2) // 50.9078??N, 004.5084??E
func (lls LatLonSpherical)Intersection(bearing1 Degrees, ll2 LatLon, bearing2 Degrees) LatLon {
	const ?? = math.Pi
	?? := math.Nextafter(1, 2) - 1
	
        // see www.edwilliams.org/avform.htm#Intersection

        ??1 := lls.ll.Latitude.Radians()
	??1 := lls.ll.Longitude.Radians()
        ??2 := ll2.Latitude.Radians()
	??2 := ll2.Longitude.Radians()
        ??13 := bearing1.Radians()
	??23 := bearing2.Radians()
        ???? := ??2 - ??1
	???? := ??2 - ??1

        // angular distance p1-p2
        ??12 := 2 * math.Asin(math.Sqrt(math.Sin(????/2) * math.Sin(????/2) +
		math.Cos(??1) * math.Cos(??2) * math.Sin(????/2) * math.Sin(????/2)))
        if math.Abs(??12) < ?? {
		return lls.ll  // coincident points
	}

        // initial/final bearings between points
        cos??a := (math.Sin(??2) - math.Sin(??1)*math.Cos(??12)) / (math.Sin(??12)*math.Cos(??1))
        cos??b := (math.Sin(??1) - math.Sin(??2)*math.Cos(??12)) / (math.Sin(??12)*math.Cos(??2))
        ??a := math.Acos(math.Min(math.Max(cos??a, -1), 1))       // protect against rounding errors
        ??b := math.Acos(math.Min(math.Max(cos??b, -1), 1))       // protect against rounding errors

        ??12 := ??a
	if math.Sin(??2-??1) <= 0 {
		??12 = 2*??-??a
	}
        ??21 := ??b
	if math.Sin(??2-??1)>0 {
		??21 = 2*??-??b
	}

        ??1 := ??13 - ??12    // angle 2-1-3
        ??2 := ??21 - ??23    // angle 1-2-3

        if math.Sin(??1) == 0 && math.Sin(??2) == 0 {
		return LatLon{Latitude: Degrees(math.NaN()), Longitude: Degrees(math.NaN())}  // infinite intersections
	}
        if math.Sin(??1) * math.Sin(??2) < 0 {
		return LatLon{Latitude: Degrees(math.NaN()), Longitude: Degrees(math.NaN())}  // ambiguous intersection (antipodal?)
	}

        cos??3 := -math.Cos(??1)*math.Cos(??2) + math.Sin(??1)*math.Sin(??2)*math.Cos(??12)

        ??13 := math.Atan2(math.Sin(??12)*math.Sin(??1)*math.Sin(??2), math.Cos(??2) + math.Cos(??1)*cos??3)

        ??3 := math.Asin(math.Min(math.Max(math.Sin(??1)*math.Cos(??13) + math.Cos(??1)*math.Sin(??13)*math.Cos(??13), -1), 1))

        ????13 := math.Atan2(math.Sin(??13)*math.Sin(??13)*math.Cos(??1), math.Cos(??13) - math.Sin(??1)*math.Sin(??3))
        ??3 := ??1 + ????13

        lat := DegreesFromRadians(??3)
        lon := DegreesFromRadians(??3)

        return LatLon{Latitude: Wrap90(lat), Longitude: Wrap180(lon)}
}

