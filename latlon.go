package geod

// Pure Go re-implementation of https://github.com/chrisveness/geodesy

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

import (
	"math"
	"fmt"
	"strings"
)

// Degrees angle
// Defining it as a type makes it harder to mix Degrees and Radians in your code, you're welcome :)
type Degrees float64

// Valid returns true if the angle is valid. Invalid angles are returned by
// functions when the result cannot be calculated.
func (d Degrees)Valid() bool {
	return !math.IsNaN(float64(d))
}

// ToRadians takes an argument in degrees and returns it in radians
func (d Degrees)Radians() float64 {
	return float64(d) * math.Pi / 180.0
}

// RoundTo returns the degrees as a float rounded to `n` decimal points.
func (d Degrees)RoundTo(n int) float64 {
	p10 := math.Pow10(n)
	return math.Round(p10 * float64(d)) / p10
}

// ToDegrees takes an argument in radians and returns it in degrees
func DegreesFromRadians(radians float64) Degrees {
	return Degrees(radians * 180.0 / math.Pi)
}


// LatLon represents a point on Earth defined by its Latitude and Longitude
type LatLon struct {
	Latitude Degrees
	Longitude Degrees
}

// Valid returns true if the coordinates are valid. Invalid coordinates are returned by
// functions when the result cannot be calculated.
func (ll LatLon)Valid() bool {
	if math.IsNaN(float64(ll.Latitude)) || math.IsNaN(float64(ll.Longitude)) {
		return false
	}

	return true
}

// Equals returns true if `ll` and `other` have identical Latitude and Longitude values
func (ll LatLon)Equals(other LatLon) bool {
	epsilon := math.Nextafter(1, 2) - 1
	
        if math.Abs(float64(ll.Latitude) - float64(other.Latitude)) > epsilon {
		return false
	}

        if math.Abs(float64(ll.Longitude) - float64(other.Longitude)) > epsilon {
		return false
	}

        return true
}

// ParseLatLon parses a latitude/longitude point from a variety of formats.
//
// Latitude & longitude (in degrees) can be supplied as two separate string parameters or
// as a single comma-separated lat/lon string
//
// The latitude/longitude values may be signed decimal or deg-min-sec (hexagesimal) suffixed by compass direction (NSEW)
// a variety of separators are accepted. Examples: -3.62, '3 37 12W', '3°37′12″W'.
//
// Thousands/decimal separators must be comma/dot
//
// Arguments:
// lat|latlon - Latitude (in degrees), or comma-separated lat/lon
// [lon]      - Longitude (in degrees).
//
// Returns Latitude/longitude point on WGS84 (LatLon)
//
// Example:
// p1 := ParseLatLon(51.47788, -0.00147)         // numeric pair
// p2 := ParseLatLon("51.47788", "-0.00147")     // string pair
// p3 := ParseLatLon("51°28′40″N, 000°00′05″W")   // single dms string
// p4 := ParseLatLon("51°28′40″N", "000°00′05″W") // dms lat string, dms lon string
func ParseLatLon(args ...interface{}) (LatLon, error) {
	if len(args) == 0 {
		return LatLon{}, fmt.Errorf("Invalid (empty) point")
	}

	// split the arguments into lat, lon
	var args2 []interface{}
	if len(args) == 1 {
		// single string of "lat, lon"
		s, ok := args[0].(string)
		if !ok {
			return LatLon{}, fmt.Errorf("Invalid argument type: %T", args[0])
		}
		tokens := strings.Split(s, ",")
		if len(tokens) > 2 {
			return LatLon{}, fmt.Errorf("Failed to parse argument: too many items")
		}
		if len(tokens) == 1 {
			return LatLon{}, fmt.Errorf("Failed to parse argument: latitude and longitude are required")
		}
		args2 = []interface{}{tokens[0], tokens[1]}
	} else if len(args) == 2 {
		args2 = args
	} else {
		return LatLon{}, fmt.Errorf("Too many arguments")
	}
	var lat, lon Degrees

	// we now have 2 values in args2: lat, lon
	switch v := args2[0].(type) {
	case string:
		lat0, err := ParseDMS(v)
		if err != nil {
			return LatLon{}, fmt.Errorf("Failed to parse latitude: %v", err)
		}
		lat = lat0
	case float64:
		if math.IsNaN(v) {
			return LatLon{}, fmt.Errorf("Latitude cannot be NaN")
		}
		lat = Degrees(v)
	case float32:
		if math.IsNaN(float64(v)) {
			return LatLon{}, fmt.Errorf("Latitude cannot be NaN")
		}
		lat = Degrees(v)
	case Degrees:
		lat = v
	default:
		return LatLon{}, fmt.Errorf("Invalid type for latitude: %T", v)
	}

	lat = Wrap90(lat)

	switch v := args2[1].(type) {
	case string:
		lon0, err := ParseDMS(v)
		if err != nil {
			return LatLon{}, fmt.Errorf("Failed to parse longitude: %v", err)
		}
		lon = lon0
	case float64:
		if math.IsNaN(v) {
			return LatLon{}, fmt.Errorf("Longitude cannot be NaN")
		}
		lon = Degrees(v)
	case float32:
		if math.IsNaN(float64(v)) {
			return LatLon{}, fmt.Errorf("Longitude cannot be NaN")
		}
		lon = Degrees(v)
	case Degrees:
		lon = v
	default:
		return LatLon{}, fmt.Errorf("Invalid type for longitude: %T", v)
	}

	lon = Wrap180(lon)

        return LatLon{Latitude: lat, Longitude: lon}, nil
}
