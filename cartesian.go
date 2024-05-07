package geod

// Pure Go re-implementation of https://github.com/chrisveness/geodesy

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

import (
	"math"
)

// Cartesian represents ECEF (earth-centered earth-fixed) geocentric cartesian coordinates
type Cartesian Vector3D

// Cartesian converts the point from (geodetic) latitude/longitude coordinates to (geocentric) cartesian (x/y/z) coordinates
// Returns the Cartesian point equivalent to lat/lon point, with x, y, z in metres from earth centre.
func (l LatLonEllipsoidal) Cartesian() Cartesian {
	// x = (ν+h)⋅cosφ⋅cosλ, y = (ν+h)⋅cosφ⋅sinλ, z = (ν⋅(1-e²)+h)⋅sinφ
	// where ν = a/√(1−e²⋅sinφ⋅sinφ), e² = (a²-b²)/a² or (better conditioned) 2⋅f-f²
	ellipsoid := l.ellipsoid

	φ := l.Latitude.Radians()
	λ := l.Longitude.Radians()
	h := l.Height
	a := ellipsoid.a
	f := ellipsoid.f

	sinφ := math.Sin(φ)
	cosφ := math.Cos(φ)
	sinλ := math.Sin(λ)
	cosλ := math.Cos(λ)

	eSq := 2*f - f*f                    // 1st eccentricity squared ≡ (a²-b²)/a²
	ν := a / math.Sqrt(1-eSq*sinφ*sinφ) // radius of curvature in prime vertical

	x := (ν + h) * cosφ * cosλ
	y := (ν + h) * cosφ * sinλ
	z := (ν*(1-eSq) + h) * sinφ

	return Cartesian{x, y, z}
}

// LatLonEllipsoidal converts this (geocentric) cartesian (x/y/z) coordinate to a (geodetic) latitude/longitude point
// on specified ellipsoid.
// Uses Bowring’s (1985) formulation for μm precision in concise form; `The accuracy of geodetic latitude and height equations'
// B R Bowring, Survey Review vol 28, 218, Oct 1985.
//
// Argument
//
//	ellipsoid - the Ellipsoid to use for the conversion
//
// Returns LatLonEllipsoidal - Latitude/longitude point defined by cartesian coordinates, on given ellipsoid.
//
// Example
// c := geod.Cartesian{X: 4027893.924, Y: 307041.993, Z: 4919474.294}
// p := c.LatLon(geod.WGS84())   // 50.7978°N, 004.3592°E
func (c Cartesian) LatLonEllipsoidal(ellipsoid Ellipsoid) LatLonEllipsoidal {
	// note ellipsoid is available as a parameter for when LatLon is used in EllipsoidalDatum / EllipsoidalReferenceframe.

	x := c.X
	y := c.Y
	z := c.Z
	a := ellipsoid.a
	b := ellipsoid.b
	f := ellipsoid.f

	e2 := 2*f - f*f           // 1st eccentricity squared ≡ (a²−b²)/a²
	ε2 := e2 / (1 - e2)       // 2nd eccentricity squared ≡ (a²−b²)/b²
	p := math.Sqrt(x*x + y*y) // distance from minor axis
	R := math.Sqrt(p*p + z*z) // polar radius

	// parametric latitude (Bowring eqn.17, replacing tanβ = z·a / p·b)
	tanβ := (b * z) / (a * p) * (1 + ε2*b/R)
	sinβ := tanβ / math.Sqrt(1+tanβ*tanβ)
	cosβ := sinβ / tanβ

	// geodetic latitude (Bowring eqn.18: tanφ = z+ε²⋅b⋅sin³β / p−e²⋅cos³β)
	var φ float64
	if !math.IsNaN(cosβ) {
		φ = math.Atan2(z+ε2*b*sinβ*sinβ*sinβ, p-e2*a*cosβ*cosβ*cosβ)
	}

	// longitude
	λ := math.Atan2(y, x)

	// height above ellipsoid (Bowring eqn.7)
	sinφ := math.Sin(φ)
	cosφ := math.Cos(φ)
	ν := a / math.Sqrt(1-e2*sinφ*sinφ) // length of the normal terminated by the minor axis
	h := p*cosφ + z*sinφ - (a * a / ν)

	return LatLonEllipsoidal{
		LatLon: LatLon{
			Latitude:  DegreesFromRadians(φ),
			Longitude: DegreesFromRadians(λ),
		},
		Height:    h,
		ellipsoid: ellipsoid,
	}
}
