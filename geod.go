package geod

// Pure Go re-implementation of https://github.com/chrisveness/geodesy

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

// Model defines the Earth model used for calculations.
// The following models are implemented:
//    geod.SphericalModel  - spherical Earth, along great circles
//    geod.RhumbModel      - spherical Earth, along rhumb lines
//    geod.VincentyModel   - ellipsoid Earth, high accuracy, slower than SphericalModel
type Model interface {
	DistanceTo(ll LatLon) DistanceUnits
	InitialBearingTo(ll LatLon) Degrees
	FinalBearingOn(ll LatLon) Degrees
	DestinationPoint(distance float64, bearing Degrees) LatLon
	MidPointTo(ll LatLon) LatLon
	IntermediatePointTo(ll LatLon, fraction float64) LatLon
	IntermediatePointsTo(ll LatLon, fractions []float64) []LatLon
	LatLon() LatLon
}

type EarthModel func(LatLon, ...interface{}) Model

// MidPoint returns the point halfway between `start` and `end` using the given `model`.
//
// Arguments:
//
// start - starting point
// end - end point (destination)
// model - a function that converts a `LatLon` to a structure appropriate for the `Model` to be used
//         This is how you select the model you wish to use for the calculations. See the description of `Model`
//         for list of available functions.
// modelArgs - additional arguments to pass to the `model` function, if needed, for example the `Ellipsoid`
//         for ellipsoid models.
//
// Returns the halfway point.
// If the point cannot be calculated an invalid point is returned, which can be tested using `LatLon.Valid()`
//
// Example:
// p1 := geod.NewLatLon(10.1, -20.0)
// p2 := geod.NewLatLon(12.1, -23.2)
// mid := geod.MidPoint(p1, p2, geod.SphericalModel)
func MidPoint(start, end LatLon, model EarthModel, modelArgs ...interface{}) LatLon {
	p1 := model(start, modelArgs...)
	return p1.MidPointTo(end)
}

// Distance returns the distance in `DistanceUnits` between points `start` and `end` using the given `model`.
//
// Arguments:
//
// start - starting point
// end - end point (destination)
// model - a function that converts a `LatLon` to a structure appropriate for the `Model` to be used
//         This is how you select the model you wish to use for the calculations. See the description of `Model`
//         for list of available functions.
// modelArgs - additional arguments to pass to the `model` function, if needed, for example the `Ellipsoid`
//         for ellipsoid models.
//
// Returns the distance in `DistanceUnits`
// If the distance cannot be calculated an invalid  is returned, which can be tested using `DistanceUnits.Valid()`
//
// Example:
// p1 := geod.NewLatLon(10.1, -20.0)
// p2 := geod.NewLatLon(12.1, -23.2)
// dist := geod.MidPoint(p1, p2, geod.VincentyModel, WGS84)    // WGS84 can be omitted, it's the default and only
//                                                                `Ellipsoid` currently defined
// metres := dist.Metres()
func Distance(start, end LatLon, model EarthModel, modelArgs ...interface{}) DistanceUnits {
	p1 := model(start, modelArgs...)
	return p1.DistanceTo(end)
}

// InitialBearing returns the initial bearing going from `start` to `end` using the given `model`.
//
// Arguments:
//
// start - starting point
// end - end point (destination)
// model - a function that converts a `LatLon` to a structure appropriate for the `Model` to be used
//         This is how you select the model you wish to use for the calculations. See the description of `Model`
//         for list of available functions.
// modelArgs - additional arguments to pass to the `model` function, if needed, for example the `Ellipsoid`
//         for ellipsoid models.
//
// Returns the initial bearing in `Degrees` from North
// If the bearing cannot be calculated NaN value is returned, which can be tested using `Degrees.Valid()`
//
// Example:
// p1 := geod.NewLatLon(10.1, -20.0)
// p2 := geod.NewLatLon(12.1, -23.2)
// bearing := geod.InitialBearing(p1, p2, geod.SphericalModel)
func InitialBearing(start, end LatLon, model EarthModel, modelArgs ...interface{}) Degrees {
	p1 := model(start, modelArgs...)
	return p1.InitialBearingTo(end)
}

// FinalBearing returns the final bearing having travelled from `start` to `end` using the given `model`.
//
// Arguments:
//
// start - starting point
// end - end point (destination)
// model - a function that converts a `LatLon` to a structure appropriate for the `Model` to be used
//         This is how you select the model you wish to use for the calculations. See the description of `Model`
//         for list of available functions.
// modelArgs - additional arguments to pass to the `model` function, if needed, for example the `Ellipsoid`
//         for ellipsoid models.
//
// Returns the final bearing in `Degrees` from North
// If the bearing cannot be calculated NaN value is returned, which can be tested using `Degrees.Valid()`
//
// Example:
// p1 := geod.NewLatLon(10.1, -20.0)
// p2 := geod.NewLatLon(12.1, -23.2)
// bearing := geod.FinalBearing(p1, p2, geod.SphericalModel)
func FinalBearing(start, end LatLon, model EarthModel, modelArgs ...interface{}) Degrees {
	p1 := model(start, modelArgs...)
	return p1.FinalBearingOn(end)
}

// DestinationPoint returns the destination point going from `start` having travelled `distance` on the given initial bearing,
// using the given `model`.
//
// Arguments:
//
// start - starting point
// distance - distance travelled, in metres -- Note: I might change this to DistanceUnits in the future (FIXME)
// bearing - initial bearing in `Degrees` from North
// model - a function that converts a `LatLon` to a structure appropriate for the `Model` to be used
//         This is how you select the model you wish to use for the calculations. See the description of `Model`
//         for list of available functions.
// modelArgs - additional arguments to pass to the `model` function, if needed, for example the `Ellipsoid`
//         for ellipsoid models.
//
// Returns the final point (destination)
// If the point cannot be calculated an invalid point is returned, which can be tested using `LatLon.Valid()`
//
// Example:
// p1 := geod.NewLatLon(10.1, -20.0)
// bearing := geod.Degrees(23.2)
// p2 := geod.Destination(p1, 100000.0, bearing, geod.RhumbModel) // 100kms from p1 heading 23.2 along a rhumb line
func DestinationPoint(start LatLon, distance float64, bearing Degrees, model EarthModel,
	modelArgs ...interface{}) LatLon {

	p1 := model(start, modelArgs...)
	return p1.DestinationPoint(distance, bearing)
}

// IntermediatePoint returns the point at the given fraction between `start` and `end`.
//
// Arguments:
//
// start - starting point
// end - end point (destination)
// fraction - the fraction between the two points (0.0 = `start`, 1.0 = `end`)
// model - a function that converts a `LatLon` to a structure appropriate for the `Model` to be used
//         This is how you select the model you wish to use for the calculations. See the description of `Model`
//         for list of available functions.
// modelArgs - additional arguments to pass to the `model` function, if needed, for example the `Ellipsoid`
//         for ellipsoid models.
//
// Returns the intermediate point at the given fraction.
// If the point cannot be calculated an invalid point is returned, which can be tested using `LatLon.Valid()`
//
// Example:
// p1 := geod.NewLatLon(10.1, -20.0)
// p2 := geod.NewLatLon(12.1, -23.2)
// pInt := geod.IntermediatePoint(p1, p2, 0.24, geod.VincentyModel)
func IntermediatePoint(start, end LatLon, fraction float64, model EarthModel,
	modelArgs ...interface{}) LatLon {

	p1 := model(start, modelArgs...)
	return p1.IntermediatePointTo(end, fraction)
}

// IntermediatePoints returns a slice of points at the given fractions between `start` and `end`.
// This is far more efficient than multiple `IntermediatePoint` calls in a loop as some of the
// expensive calculations are reused and each franctional point is calculated in parallel.
//
// Arguments:
//
// start - starting point
// end - end point (destination)
// fractions - slice of fractions between the two points (0.0 = `start`, 1.0 = `end`)
// model - a function that converts a `LatLon` to a structure appropriate for the `Model` to be used
//         This is how you select the model you wish to use for the calculations. See the description of `Model`
//         for list of available functions.
// modelArgs - additional arguments to pass to the `model` function, if needed, for example the `Ellipsoid`
//         for ellipsoid models.
//
// Returns slice of intermediate points at the given fractions.
// Points that cannot be calculated are returned as invalid points, can be tested using `LatLon.Valid()`
//
// Example:
// p1 := geod.NewLatLon(10.1, -20.0)
// p2 := geod.NewLatLon(12.1, -23.2)
// pInt := geod.IntermediatePoint(p1, p2, 0.24, geod.VincentyModel)
func IntermediatePoints(start, end LatLon, fractions []float64, model EarthModel,
	modelArgs ...interface{}) []LatLon {

	p1 := model(start, modelArgs...)
	return p1.IntermediatePointsTo(end, fractions)
}
