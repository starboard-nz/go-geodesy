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
 * Distances & bearings between points, and destination points given start points & initial bearings,
 * calculated on an ellipsoidal earth model using ‘direct and inverse solutions of geodesics on the
 * ellipsoid’ devised by Thaddeus Vincenty.
 *
 * From: T Vincenty, "Direct and Inverse Solutions of Geodesics on the Ellipsoid with application of
 * nested equations", Survey Review, vol XXIII no 176, 1975. www.ngs.noaa.gov/PUBS_LIB/inverse.pdf.
 */

import (
	"math"
	"sync"
)

// LatLonEllipsoidalVincenty represents a point used for calculations using a the Vincenty method, on an
// ellipsoidal Earth model.
type LatLonEllipsoidalVincenty struct {
	ll LatLon
	ellipsoid Ellipsoid
}

// VincentyModel returns a `Model` that wraps geodesy calculations using the Vincenty method on an ellipsoidal Earth model
func VincentyModel(ll LatLon, modelArgs ...interface{}) Model {
	ellipsoid := WGS84()
	if len(modelArgs) != 0 {
		if len(modelArgs) > 1 {
			panic("Invalid number of arguments in call to VincentyModel()")
		}
		switch v := modelArgs[0].(type) {
		case Ellipsoid:
			ellipsoid = v
		case func() Ellipsoid:
			ellipsoid = v()
		default:
			panic("Invalid argument type in call to VincentyModel()")
		}
	}
	return LatLonEllipsoidalVincenty{ll: ll, ellipsoid: ellipsoid}
}

// LatLon converts LatLonEllipsoidalVincenty to LatLon
func (llv LatLonEllipsoidalVincenty)LatLon() LatLon {
	return llv.ll
}

// NewLatLonEllipsodialVincenty creates a new LatLonEllipsoidalVincenty struct
func NewLatLonEllipsodialVincenty(latitude, longitude float64, ellipsoid Ellipsoid) LatLonEllipsoidalVincenty {
	return LatLonEllipsoidalVincenty{
		ll: LatLon{
			Latitude: Wrap90(Degrees(latitude)),
			Longitude: Wrap180(Degrees(longitude)),
		},
		ellipsoid: ellipsoid,
	}
}

// VincentyDirect - Vincenty direct calculation - calculates the destination point and final bearing given the
// starting point, distance and initial bearing.
//
// Arguments
//
// distance - Distance along bearing in metres
// initialBearing - Initial bearing in degrees from North
//
// Returns (destination, finalBearing)
func (llv LatLonEllipsoidalVincenty)VincentyDirect(distance float64, initialBearing Degrees) (LatLon, Degrees) {
        φ1 := llv.ll.Latitude.Radians()
	λ1 := llv.ll.Longitude.Radians()
        α1 := initialBearing.Radians()
        s := distance

	a := llv.ellipsoid.a
	b := llv.ellipsoid.b
	f := llv.ellipsoid.f

        sinα1 := math.Sin(α1)
        cosα1 := math.Cos(α1)

        tanU1 := (1-f) * math.Tan(φ1)
	cosU1 := 1 / math.Sqrt((1 + tanU1*tanU1))
	sinU1 := tanU1 * cosU1

        σ1 := math.Atan2(tanU1, cosα1)    // σ1 = angular distance on the sphere from the equator to P1
        sinα := cosU1 * sinα1             // α = azimuth of the geodesic at the equator
        cosSqα := 1 - sinα*sinα
        uSq := cosSqα * (a*a - b*b) / (b*b)
        A := 1 + uSq/16384*(4096+uSq*(-768+uSq*(320-175*uSq)))
        B := uSq/1024 * (256+uSq*(-128+uSq*(74-47*uSq)))

        σ := s / (b*A)

	var sinσ, cosσ float64
	var Δσ float64                    // σ = angular distance P₁ P₂ on the sphere
        var cos2σₘ float64                // σₘ = angular distance on the sphere from the equator to the midpoint of the line

        var σʹ float64
	iterations := 0
        for {
		cos2σₘ = math.Cos(2*σ1 + σ)
		sinσ = math.Sin(σ)
		cosσ = math.Cos(σ)
		Δσ = B * sinσ * (cos2σₘ + B/4 * (cosσ * (-1 + 2 * cos2σₘ * cos2σₘ) -
			B/6 * cos2σₘ * (-3 + 4 * sinσ * sinσ) * (-3 + 4 * cos2σₘ * cos2σₘ)))
		σʹ = σ
		σ = s / (b * A) + Δσ
		iterations++
		if math.Abs(σ - σʹ) <= 1e-12 || iterations >= 100 {
			break
		}
	}
	if iterations >= 100 {
		// algorithm failed to converge
		return LatLon{Latitude: Degrees(math.NaN()), Longitude: Degrees(math.NaN())}, Degrees(math.NaN())
	}

        x := sinU1 * sinσ - cosU1 * cosσ * cosα1
        φ2 := math.Atan2(sinU1 * cosσ + cosU1 *sinσ *cosα1, (1 - f) * math.Sqrt(sinα * sinα + x*x))
        λ := math.Atan2(sinσ * sinα1, cosU1 * cosσ - sinU1 * sinσ * cosα1)
        C := f / 16 * cosSqα * (4 + f * (4 - 3 * cosSqα))
        L := λ - (1-C) * f * sinα * (σ + C * sinσ * (cos2σₘ + C * cosσ * (-1 + 2 * cos2σₘ * cos2σₘ)))
        λ2 := λ1 + L

        α2 := math.Atan2(sinα, -x)

        destinationPoint := LatLon{Latitude: Wrap90(DegreesFromRadians(φ2)), Longitude: Wrap180(DegreesFromRadians(λ2))}
	finalBearing := Wrap360(DegreesFromRadians(α2))

        return destinationPoint, finalBearing
}

// VincentyInverse - Vincenty inverse calculation.  Calculates the distance, initial and final bearing going
// from point `llv` to `dest`, using the Vincenty method.
//
// Arguments:
//
// dest - destination point
//
// Returns (distance from `llv` to `dest`, initial bearing in degrees from North, final bearing in degrees from North)
func (llv LatLonEllipsoidalVincenty)VincentyInverse(dest LatLon) (units.Distance, Degrees, Degrees) {
	if llv.ll.Equals(dest) {
		return units.Metre(math.NaN()), Degrees(math.NaN()), Degrees(math.NaN())
	}

	const π = math.Pi
	ε := math.Nextafter(1, 2) - 1

        φ1 := llv.ll.Latitude.Radians()
	λ1 := llv.ll.Longitude.Radians()
        φ2 := dest.Latitude.Radians()
	λ2 := dest.Longitude.Radians()

	a := llv.ellipsoid.a
	b := llv.ellipsoid.b
	f := llv.ellipsoid.f

        L := λ2 - λ1   // L = difference in longitude, U = reduced latitude, defined by tan U = (1-f)·tanφ.
        tanU1 := (1.0 - f) * math.Tan(φ1)
	cosU1 := 1.0 / math.Sqrt((1 + tanU1 * tanU1))
	sinU1 := tanU1 * cosU1
	
        tanU2 := (1.0 - f) * math.Tan(φ2)
	cosU2 := 1 / math.Sqrt((1 + tanU2 * tanU2))
	sinU2 := tanU2 * cosU2

        isAntipodal := math.Abs(L) > π/2 || math.Abs(φ2 - φ1) > π/2

        λ := L
	var sinλ, cosλ float64          // λ = difference in longitude on an auxiliary sphere
	var sinSqσ float64              // σ = angular distance P₁ P₂ on the sphere
        σ := 0.0
	sinσ := 0.0
	cosσ := 1.0
	if isAntipodal {
		σ = π
		cosσ = -1.0
	}
        cos2σₘ := 1.0                   // σₘ = angular distance on the sphere from the equator to the midpoint of the line
        var sinα float64                // α = azimuth of the geodesic at the equator
	cosSqα := 1.0

        var C, λʹ, iterationCheck float64
	iterations := 0
        for {
		sinλ = math.Sin(λ)
		cosλ = math.Cos(λ)
		sinSqσ = (cosU2 * sinλ) * (cosU2 * sinλ) + (cosU1 * sinU2 - sinU1 * cosU2 * cosλ) *
			(cosU1 * sinU2 - sinU1 * cosU2 * cosλ)
		if math.Abs(sinSqσ) < ε {
			break           // co-incident/antipodal points (falls back on λ/σ = L)
		}
		sinσ = math.Sqrt(sinSqσ)
		cosσ = sinU1 * sinU2 + cosU1 * cosU2 * cosλ
		σ = math.Atan2(sinσ, cosσ)
		sinα = cosU1 * cosU2 * sinλ / sinσ
		cosSqα = 1 - sinα * sinα
		if cosSqα != 0 {
			cos2σₘ = cosσ - 2 * sinU1 * sinU2 / cosSqα
		} else {
			cos2σₘ = 0.0     // on equatorial line cos²α = 0 (§6)
		}
		C = f / 16 * cosSqα * (4 + f * (4 - 3 * cosSqα))
		λʹ = λ
		λ = L + (1 - C) * f * sinα * (σ + C * sinσ * (cos2σₘ + C * cosσ * (-1 + 2 * cos2σₘ * cos2σₘ)))
		if isAntipodal {
			iterationCheck = math.Abs(λ) - π
		} else {
			iterationCheck = math.Abs(λ)
		}
		if (iterationCheck > π) {
			return units.Metre(math.NaN()), Degrees(math.NaN()), Degrees(math.NaN())
		}
		iterations++
		if math.Abs(λ - λʹ) <= 1e-12 || iterations >= 1000 {
			break
		}
        }

        if iterations >= 1000 {
		return units.Metre(math.NaN()), Degrees(math.NaN()), Degrees(math.NaN())
	}

        uSq := cosSqα * (a * a - b * b) / (b * b)
        A := 1 + uSq / 16384 * (4096 + uSq * (-768 + uSq * (320 - 175 * uSq)))
        B := uSq / 1024 * (256 + uSq * (-128 + uSq * (74 - 47 * uSq)))
        Δσ := B * sinσ * (cos2σₘ + B / 4 * (cosσ * (-1 + 2 * cos2σₘ * cos2σₘ) -
		B / 6 * cos2σₘ * (-3 + 4 * sinσ * sinσ) * (-3 + 4 * cos2σₘ * cos2σₘ)))

        s := b * A * (σ - Δσ)      // s = length of the geodesic

        // note special handling of exactly antipodal points where sin²σ = 0 (due to discontinuity
        // atan2(0, 0) = 0 but atan2(ε, 0) = π/2 / 90°) - in which case bearing is always meridional,
        // due north (or due south!)
        // α = azimuths of the geodesic; α2 the direction P₁ P₂ produced
        α1 := 0.0
	if math.Abs(sinSqσ) >= ε {
		α1 = math.Atan2(cosU2*sinλ, cosU1 * sinU2 - sinU1 * cosU2 * cosλ)
	}
	α2 := π
        if math.Abs(sinSqσ) >= ε {
		α2 = math.Atan2(cosU1 * sinλ, -sinU1 * cosU2 + cosU1 * sinU2 * cosλ)
	}
	initialBearing := Degrees(math.NaN())
	if math.Abs(s) >= ε {
		initialBearing = Wrap360(DegreesFromRadians(α1))
	}
	finalBearing := Degrees(math.NaN())
	if math.Abs(s) >= ε {
		finalBearing = Wrap360(DegreesFromRadians(α2))
	}
        return units.Metre(s), initialBearing, finalBearing
}

// DistanceTo returns the distance along the surface of the earth from `llv` to `dest` using Vincenty Inverse calculation
//
// Argument:
//
// dest  - destination point
//
// Returns the `Distance` between this point and destination point in DistanceUnits
//
// Examples:
// p1 := geod.NewLatLonEllipsodialVincenty(52.205, 0.119, geod.WGS84())
// p2 := geod.LatLon{48.857, 2.351}
// d := p1.DistanceTo(p2).Metre()       // 404.3×10³ m
// m := p1.DistanceTo(p2, 3959).Mile()  // 251.2 miles
func (llv LatLonEllipsoidalVincenty)DistanceTo(dest LatLon) units.Distance {
	dist, _, _ := llv.VincentyInverse(dest)
	return dist
}

// InitialBearingTo returns the initial bearing (forward azimuth) to travel along a geodesic from `llv` to `dest`
// using the Vincenty inverse solution
//
// Arguments:
//
// dest - destination point
//
// Returns the initial bearing in degrees from North (0°..360°) or NaN if failed to converge
//
// Example:
// p1 := geod.NewLatLonEllipsodialVincenty(50.06632, -5.71475, geod.WGS84())
// p2 := geod.LatLon{58.64402, -3.07009}
// b1 := p1.InitialBearingTo(p2)    // 9.1419°
func (llv LatLonEllipsoidalVincenty)InitialBearingTo(dest LatLon) Degrees {
	_, initialBearing, _ := llv.VincentyInverse(dest)
	return initialBearing
}

// FinalBearingOn returns the final bearing (review azimuth) having travelled along a geodesic from `llv` to `dest`
// using the Vincenty inverse solution
//
// Arguments:
//
// dest - destination point
//
// Returns the final bearing in degrees from North (0°..360°) or NaN if failed to converge
//
// Example:
// p1 := geod.NewLatLonEllipsodialVincenty(50.06632, -5.71475, geod.WGS84())
// p2 := geod.LatLon{58.64402, -3.07009}
// b1 := p1.FinalBearingOn(p2)    // 11.2972°
func (llv LatLonEllipsoidalVincenty)FinalBearingOn(dest LatLon) Degrees {
	_, _, finalBearing := llv.VincentyInverse(dest)
	return finalBearing
}

// MidPointTo returns the midpoint between `llv` and `dest`.
//
// Argument:
//
// dest  - destination point
//
// Returns the middle point
//
// Example:
// p1 := geod.NewLatLonEllipsodialVincenty(52.205, 0.119, geod.WGS84())
// p2 := geod.LatLon{48.857, 2.351}
// pMid := p1.MidPointTo(p2)
func (llv LatLonEllipsoidalVincenty)MidPointTo(dest LatLon) LatLon {
	distance, initialBearing, _ := llv.VincentyInverse(dest)
	point, _ := llv.VincentyDirect(float64(distance.Metre() / 2), initialBearing)
	return point
}

// IntermediatePointsTo returns the points at the given fractions between `llv` and `dest`.
//
// Arguments:
//
// dest  - destination point
// fraction - Slice of fractions between the two points (0 = `llv`, 1 = `dest`)
//
// Returns an intermediate point for each fraction
//
// Example:
// p1 := geod.NewLatLonEllipsodialVincenty(52.205, 0.119, geod.WGS84())
// p2 := geod.LatLon{48.857, 2.351}
// pInt := p1.IntermediatePointsTo(p2, []float64{0.25, 0.5, 0.75})
func (llv LatLonEllipsoidalVincenty)IntermediatePointsTo(dest LatLon, fractions []float64) []LatLon {
	waitGroup := &sync.WaitGroup{}

	distance, initialBearing, _ := llv.VincentyInverse(dest)

	points := make([]LatLon, len(fractions))
	for i, fraction := range(fractions) {
		waitGroup.Add(1)
		go func(i int, fraction float64) {
			points[i], _ = llv.VincentyDirect(float64(distance.Metre()) * fraction, initialBearing)
			waitGroup.Done()
		} (i, fraction)
	}

	// wait for all goroutines to finish
	waitGroup.Wait()

	return points
}

// IntermediatePointTo returns the points at the given fraction between `llv` and `dest`.
//
// Arguments:
//
// dest  - destination point
// fraction - Fractions between the two points (0 = `llv`, 1 = `dest`)
//
// Returns the intermediate point.
//
// Example:
// p1 := geod.NewLatLonEllipsodialVincenty(52.205, 0.119, geod.WGS84())
// p2 := geod.LatLon{48.857, 2.351}
// pInt := p1.IntermediatePointTo(p2, 0.25)
func (llv LatLonEllipsoidalVincenty)IntermediatePointTo(dest LatLon, fraction float64) LatLon {
	distance, initialBearing, _ := llv.VincentyInverse(dest)

	point, _ := llv.VincentyDirect(float64(distance.Metre()) * fraction, initialBearing)
	return point
}

// DestinationPoint returns the destination point having travelled the given `distance` along a geodesic given by
// `initialBearing` from `llv`, using Vincenty direct solution
//
// Arguments:
//
// distance - Distance travelled along the geodesic in metres
// initialBearing - Initial bearing in degrees from North
//
// Returns the destination point
//
// Example
// p1 := geod.NewLatLonEllipsodialVincenty(-37.95103, 144.42487, geod.WGS84())
// p2 := p1.DestinationPoint(54972.271, geod.Degrees(306.86816))    // 37.6528°S, 143.9265°E
func (llv LatLonEllipsoidalVincenty)DestinationPoint(distance float64, bearing Degrees) LatLon {
	point, _ := llv.VincentyDirect(distance, bearing)
	return point
}
