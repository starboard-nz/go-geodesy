package geod

// Pure Go re-implementation of https://github.com/chrisveness/geodesy

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

// Model defines the Earth model used for calculations.
// Currently defined models:
//    geod.SphericalModel  - spherical Earth, along great circles
//    geod.RhumbModel      - spherical Earth, along rhumb lines
type Model interface {
	DistanceTo(ll LatLon) DistanceUnits
	InitialBearingTo(ll LatLon) Degrees
	FinalBearingTo(ll LatLon) Degrees
	DestinationPoint(distance float64, bearing Degrees) LatLon
	MidPointTo(ll LatLon) LatLon
	IntermediatePointTo(ll LatLon, fraction float64) LatLon
	LatLon() LatLon
}

// MidPoint returns the point halfway between `start` and `end` using the given `model`.
func MidPoint(start, end LatLon, model func(ll LatLon) Model) LatLon {
	p1 := model(start)
	return p1.MidPointTo(end)
}

// Distance returns the distance in `DistanceUnits` between points `start` and `end` using the given `model`.
func Distance(start, end LatLon, model func(ll LatLon) Model) DistanceUnits {
	p1 := model(start)
	return p1.DistanceTo(end)
}

// InitialBearing returns the initial bearing going from `start` to `end` using the given `model`.
func InitialBearing(start, end LatLon, model func(ll LatLon) Model) Degrees {
	p1 := model(start)
	return p1.InitialBearingTo(end)
}

// FinalBearing returns the final bearing going from `start` to `end` using the given `model`.
func FinalBearing(start, end LatLon, model func(ll LatLon) Model) Degrees {
	p1 := model(start)
	return p1.FinalBearingTo(end)
}

// DestinationPoint returns the destination point going from `start` having travelled `distance` on the given initial bearing,
// using the given `model`.
func DestinationPoint(start LatLon, distance float64, bearing Degrees, model func(ll LatLon) Model) LatLon {
	p1 := model(start)
	return p1.DestinationPoint(distance, bearing)
}

// IntermediatePoint returns the point at the given fraction between `start` and `end`.
func IntermediatePoint(start, end LatLon, fraction float64, model func(ll LatLon) Model) LatLon {
	p1 := model(start)
	return p1.IntermediatePointTo(end, fraction)
}
