package geod

// Pure Go re-implementation of https://github.com/chrisveness/geodesy

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

/**
 * Library of geodesy functions for operations on a spherical earth model along rhumb lines
 *
 * Includes distances, bearings, destinations, etc and other related functions.
 *
 * All calculations are done using simple spherical trigonometric formulae.
 */

import (
	"math"
	"sync"

	"github.com/starboard-nz/units"
)

// LatLonRhumb represents a point used for calculations using a spherical Earth model, along rhumb lines
type LatLonRhumb struct {
	ll LatLon
}

// RhumbModel returns a `Model` that wraps geodesy calculations using spherical Earth model along rhumb lines
func RhumbModel(ll LatLon, modelArgs ...interface{}) Model {
	if len(modelArgs) != 0 {
		panic("Invalid number of arguments in call to VincentyModel()")
	}
	return LatLonRhumb{ll: ll}
}

// LatLon converts LatLonRhumb to LatLon
func (llr LatLonRhumb)LatLon() LatLon {
	return llr.ll
}

// NewLatLonRhumb creates a new LatLonRhumb struct
func NewLatLonRhumb(latitude, longitude Degrees) LatLonRhumb {
	return LatLonRhumb{
		ll: LatLon{
			Latitude: Wrap90(latitude),
			Longitude: Wrap180(longitude),
		},
	}
}

// DistanceTo returns the distance along a rhumb line from `llr` to `dest`.
//
// Argument:
//
// dest  - destination point
//
// Returns the `Distance` between this point and destination point in Distance units.
//
// Examples:
// p1 := geod.NewLatLonRhumb(51.127, 1.338)
// p2 := geod.NewLatLonRhumb(50.964, 1.853)
// d := p1.DistanceTo(p2).Km()  //  40.31 km
func (llr LatLonRhumb)DistanceTo(dest LatLon) units.Distance {
        // see www.edwilliams.org/avform.htm#Rhumb

	const π = math.Pi
        R := earthRadius
        φ1 := llr.ll.Latitude.Radians()
        φ2 := dest.Latitude.Radians()
        Δφ := φ2 - φ1
        Δλ := Degrees(math.Abs(float64(dest.Longitude) - float64(llr.ll.Longitude))).Radians()
        // if Δλ over 180° take shorter rhumb line across the anti-meridian:
        if math.Abs(Δλ) > π {
		if Δλ > 0 {
			Δλ = -(2 * π - Δλ)
		} else {
			Δλ = 2 * π + Δλ
		}
	}

        // on Mercator projection, longitude distances shrink by latitude; q is the 'stretch factor'
        // q becomes ill-conditioned along E-W line (0/0); use empirical tolerance to avoid it
        Δψ := math.Log(math.Tan(φ2 / 2 + π / 4) / math.Tan(φ1 / 2 + π / 4))
	var q float64
        if math.Abs(Δψ) > 10e-12 {
		q = Δφ / Δψ
	} else {
		q = math.Cos(φ1)
	}

        // distance is pythagoras on 'stretched' Mercator projection, √(Δφ² + q²·Δλ²)
        δ := math.Sqrt(Δφ*Δφ + q*q * Δλ*Δλ)  // angular distance in radians
        d := δ * R

        return units.Metre(d)
}

// InitialBearingTo returns the bearing from `lls` to `dest`. In the case of rhumb lines the bearing is constant, so
// this is the same as the final bearing.
//
// Argument:
//
// dest  - destination point
//
// Returns the rhumb bearing in `Degrees` from North (0°..360°)
//
// Example:
// p1 := geod.NewLatLonRhumb(51.127, 1.338)
// p2 := geod.NewLatLonRhumb(50.964, 1.853)
// b1 := p1.InitialBearingTo(p2)    // 116.7°
func (llr LatLonRhumb)InitialBearingTo(dest LatLon) Degrees {
	if llr.ll.Equals(dest) {
		return Degrees(math.NaN())    // coincident points
	}

	const π = math.Pi
        φ1 := llr.ll.Latitude.Radians()
        φ2 := dest.Latitude.Radians()
        Δλ := (dest.Longitude - llr.ll.Longitude).Radians()

        // if dLon over 180° take shorter rhumb line across the anti-meridian:
        if math.Abs(Δλ) > π {
		if Δλ > 0 {
			Δλ = -(2 * π - Δλ)
		} else {
			Δλ = 2 * π + Δλ
		}
	}

        Δψ := math.Log(math.Tan(φ2 / 2 + π / 4) / math.Tan(φ1 / 2 + π / 4))

        θ := math.Atan2(Δλ, Δψ)

        bearing := DegreesFromRadians(θ)

        return Wrap360(bearing)
}

// FinalBearingOn returns the bearing from `lls` to `dest`. In the case of rhumb lines the bearing is constant, so
// this is the same as the initial bearing.
//
// Argument:
//
// dest  - destination point
//
// Returns the rhumb bearing in `Degrees` from North (0°..360°)
//
// Example:
// p1 := geod.NewLatLonRhumb(51.127, 1.338)
// p2 := geod.NewLatLonRhumb(50.964, 1.853)
// b1 := p1.FinalBearingOn(p2)    // 116.7°
func (llr LatLonRhumb)FinalBearingOn(dest LatLon) Degrees {
	return llr.InitialBearingTo(dest)
}

// DestinationPoint returns the destination point from `lls` having travelled the given distance
// along a rhumb line on the given bearing.
//
// Arguments:
//
// distance - Distance travelled in metres
// bearing - Bearing in `Degrees` from North
//
// Returns the destination point.
//
// Example:
// p1 := geod.NewLatLonRhumb(51.127, 1.338)
// p2 := p1.DestinationPoint(40300, geod.Degrees(116.7)) // 50.9642°N, 001.8530°E
func (llr LatLonRhumb)DestinationPoint(distance float64, bearing Degrees) LatLon {
	const π = math.Pi
        φ1 := llr.ll.Latitude.Radians()
	λ1 := llr.ll.Longitude.Radians()
        θ := bearing.Radians()

        δ := distance / earthRadius     // angular distance in radians

        Δφ := δ * math.Cos(θ)
        φ2 := φ1 + Δφ

        // check for some daft bugger going past the pole, normalise latitude if so
        if math.Abs(φ2) > π / 2 {
		if φ2 > 0 {
			φ2 = π - φ2
		} else {
			φ2 =-π - φ2
		}
	}

        Δψ := math.Log(math.Tan(φ2 / 2 + π / 4) / math.Tan(φ1 / 2 + π / 4))
        var q float64
	if math.Abs(Δψ) > 10e-12 {
		q = Δφ / Δψ
	} else {
		q = math.Cos(φ1)      // E-W course becomes ill-conditioned with 0/0
	}

        Δλ := δ * math.Sin(θ) / q
        λ2 := λ1 + Δλ

        lat := DegreesFromRadians(φ2)
        lon := DegreesFromRadians(λ2)

        return LatLon{Latitude: Wrap90(lat), Longitude: Wrap180(lon)}
}


// MidPointTo returns the loxodromic midpoint (along a rhumb line) between `llr` and `dest`.
//
// Argument:
//
// dest  - destination point
//
// Returns the middle point
//
// Example:
// p1 := geod.NewLatLonRhumb(51.127, 1.338)
// p2 := geod.NewLatLonRhumb(50.964, 1.853)
// pMid := p1.MidPointTo(p2)    // 51.0455°N, 001.5957°E
func (llr LatLonRhumb)MidPointTo(dest LatLon) LatLon {
	const π = math.Pi
        // see mathforum.org/kb/message.jspa?messageID=148837

        φ1 := llr.ll.Latitude.Radians()
	λ1 := llr.ll.Longitude.Radians()
        φ2 := dest.Latitude.Radians()
	λ2 := dest.Longitude.Radians()

        if math.Abs(λ2 - λ1) >= π {
		λ1 += 2 * π    // crossing anti-meridian
	}

        φ3 := (φ1 + φ2) / 2
        f1 := math.Tan(π / 4 + φ1 / 2)
        f2 := math.Tan(π / 4 + φ2 / 2)
        f3 := math.Tan(π / 4 + φ3 / 2)
        λ3 := ((λ2 - λ1) * math.Log(f3) + λ1 * math.Log(f2) - λ2 * math.Log(f1)) / math.Log(f2 / f1)

        if math.IsInf(λ3, 0) {
		λ3 = (λ1 + λ2) / 2  // parallel of latitude
	}

        lat := DegreesFromRadians(φ3)
        lon := DegreesFromRadians(λ3)

        return LatLon{Latitude: Wrap90(lat), Longitude: Wrap180(lon)}
}

// IntermediatePointTo returns the point at the given fraction between `lls` and `dest` along a rhumb line
//
// Arguments:
//
// dest  - destination point
// fraction - Fraction between the two points (0 = `lls`, 1 = `dest`)
//
// Returns the intermediate point.
//
// Example:
// p1 := geod.NewLatLonRhumb(51.127, 1.338)
// p2 := geod.NewLatLonRhumb(50.964, 1.853)
// pMid := p1.IntermediatePointTo(p2, 0.25)    // 51.08625°N, 001.46692°E
func (llr LatLonRhumb)IntermediatePointTo(dest LatLon, fraction float64) LatLon {
	if llr.ll.Equals(dest) {
		return llr.ll
	}

	dist := llr.DistanceTo(dest)
	frDist := float64(dist.Metre()) * fraction
	bearing := llr.InitialBearingTo(dest)
	return llr.DestinationPoint(frDist, bearing)
}

// IntermediatePointsTo returns the points at the given fractions between `llr` and `dest`.
//
// Arguments:
//
// dest  - destination point
// fraction - Slice of fractions between the two points (0 = `llr`, 1 = `dest`)
//
// Returns an intermediate point for each fraction
//
// Example:
// p1 := geod.NewLatLonRhumb(52.205, 0.119)
// p2 := geod.LatLon{48.857, 2.351}
// pInt := p1.IntermediatePointsTo(p2, []float64{0.25, 0.5, 0.75})
func (llr LatLonRhumb)IntermediatePointsTo(dest LatLon, fractions []float64) []LatLon {
	waitGroup := &sync.WaitGroup{}

	dist := llr.DistanceTo(dest)
	bearing := llr.InitialBearingTo(dest)

	points := make([]LatLon, len(fractions))
	for i, fraction := range(fractions) {
		waitGroup.Add(1)
		go func(i int, fraction float64) {
			frDist := float64(dist.Metre()) * fraction
			points[i] = llr.DestinationPoint(frDist, bearing)
			waitGroup.Done()
		} (i, fraction)
	}

	// wait for all goroutines to finish
	waitGroup.Wait()

	return points
}
