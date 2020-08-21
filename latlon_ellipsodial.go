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
	"strconv"
)

/**
 * A latitude/longitude point defines a geographic location on or above/below the earth's surface,
 * measured in degrees from the equator & the International Reference Meridian and in metres above
 * the ellipsoid, and based on a given datum.
 *
 * As so much modern geodesy is based on WGS-84 (as used by GPS), this module includes WGS-84
 * ellipsoid parameters, and it has methods for converting geodetic (latitude/longitude) points to/from
 * geocentric cartesian points; the latlon-ellipsoidal-datum and latlon-ellipsoidal-referenceframe
 * modules provide transformation parameters for converting between historical datums and between
 * modern reference frames.
 *
 * This module is used for both trigonometric geodesy (eg latlon-ellipsoidal-vincenty) and n-vector
 * geodesy (eg latlon-nvector-ellipsoidal), and also for UTM/MGRS mapping.
 *
 */


// LatLonEllipsoidal represents latitude/longitude points on an ellipsoidal model earth,
// with ellipsoid parameters and methods for converting points to/from cartesian (ECEF) coordinates.
//
// This is the core struct, which will usually be used via LatLonEllipsoidalDatum or
// LatLonEllipsoidalReferenceFrame.
type LatLonEllipsoidal struct {
	LatLon
	Height float64
	ellipsoid Ellipsoid
}

// NewLatLonEllipsodial creates a new LatLonEllipsoidal struct
func NewLatLonEllipsodial(latitude, longitude Degrees, height float64) LatLonEllipsoidal {
	return LatLonEllipsoidal{
		LatLon: LatLon{
			Latitude: Wrap90(latitude),
			Longitude: Wrap180(longitude),
		},
		Height: height,
		ellipsoid: WGS84(),
	}
}

// ParseLatLonEllipsoidal parses a latitude/longitude point from a variety of formats
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
// [height]   - Height above ellipsoid in metres.
//
// Returns Latitude/longitude point on WGS84 ellipsoidal model earth (LatLonEllipsoidal)
//
// Example:
// p1 := ParseLatLon(51.47788, -0.00147)         // numeric pair
// p2 := ParseLatLon("51.47788", "-0.00147")         // string pair
// p3 := ParseLatLon("51°28′40″N, 000°00′05″W", 17)   // dms string + height
// p4 := ParseLatLon("51°28′40″N", "000°00′05″W", 17) // dms lat, dms lon, height
func ParseLatLonEllipsoidal(args ...interface{}) (LatLonEllipsoidal, error) {
	if len(args) == 0 {
		return LatLonEllipsoidal{}, fmt.Errorf("Invalid (empty) point")
	}

	// split the arguments into lat, lon, height
	var args3 []interface{}
	if len(args) == 1 {
		// single string of "lat, lon[, height]"
		s, ok := args[0].(string)
		if !ok {
			return LatLonEllipsoidal{}, fmt.Errorf("Invalid argument type: %T", args[0])
		}
		tokens := strings.Split(s, ",")
		if len(tokens) > 3 {
			return LatLonEllipsoidal{}, fmt.Errorf("Failed to parse argument: too many items")
		}
		if len(tokens) == 1 {
			return LatLonEllipsoidal{}, fmt.Errorf("Failed to parse argument: at least latitude and longitude are required")
		}
		if len(tokens) == 3 {
			args3 = []interface{}{tokens[0], tokens[1], tokens[2]}
		} else {
			args3 = []interface{}{tokens[0], tokens[1], 0.0}
		}
	} else if len(args) == 2 {
		// either lat + lon or a lat/lon string + height
		s, ok := args[0].(string)
		if !ok {
			// not a string, so must be lat + lon
			args3 = []interface{}{args[0], args[1], 0.0}
		} else {
			tokens := strings.Split(s, ",")
			if len(tokens) > 2 {
				return LatLonEllipsoidal{}, fmt.Errorf("Failed to parse argument: too many items")
			}
			if len(tokens) == 1 {
				// lat + lon
				args3 = append(args, 0.0)
			} else if len(tokens) == 2 {
				// lat/lon + height
				args3 = []interface{}{tokens[0], tokens[1], args[1]}
			}
		}
	} else if len(args) == 3 {
		args3 = args
	} else {
		return LatLonEllipsoidal{}, fmt.Errorf("Too many arguments")
	}
	var lat, lon Degrees
	var height float64

	// we now have 3 values in args3: lat, lon, height
	switch v := args3[0].(type) {
	case string:
		lat0, err := ParseDMS(v)
		if err != nil {
			return LatLonEllipsoidal{}, fmt.Errorf("Failed to parse latitude: %v", err)
		}
		lat = lat0
	case float64:
		if math.IsNaN(v) {
			return LatLonEllipsoidal{}, fmt.Errorf("Latitude cannot be NaN")
		}
		lat = Degrees(v)
	case float32:
		if math.IsNaN(float64(v)) {
			return LatLonEllipsoidal{}, fmt.Errorf("Latitude cannot be NaN")
		}
		lat = Degrees(v)
	case Degrees:
		lat = v
	default:
		return LatLonEllipsoidal{}, fmt.Errorf("Invalid type for latitude: %T", v)
	}

	lat = Wrap90(lat)

	switch v := args3[1].(type) {
	case string:
		lon0, err := ParseDMS(v)
		if err != nil {
			return LatLonEllipsoidal{}, fmt.Errorf("Failed to parse longitude: %v", err)
		}
		lon = lon0
	case float64:
		if math.IsNaN(v) {
			return LatLonEllipsoidal{}, fmt.Errorf("Longitude cannot be NaN")
		}
		lon = Degrees(v)
	case float32:
		if math.IsNaN(float64(v)) {
			return LatLonEllipsoidal{}, fmt.Errorf("Longitude cannot be NaN")
		}
		lon = Degrees(v)
	case Degrees:
		lon = v
	default:
		return LatLonEllipsoidal{}, fmt.Errorf("Invalid type for longitude: %T", v)
	}

	lon = Wrap180(lon)

	switch v := args3[2].(type) {
	case string:
		var err error
		if height, err = strconv.ParseFloat(v, 64); err != nil {
			return LatLonEllipsoidal{}, fmt.Errorf("Failed to parse height: %v", err)
		}
	case float64:
		if math.IsNaN(v) {
			return LatLonEllipsoidal{}, fmt.Errorf("Height cannot be NaN")
		}
		height = v
	case float32:
		if math.IsNaN(float64(v)) {
			return LatLonEllipsoidal{}, fmt.Errorf("Height cannot be NaN")
		}
		height = float64(v)
	}
		
        return LatLonEllipsoidal{
		LatLon: LatLon{
			Latitude: lat,
			Longitude: lon,
		},
		Height: height,
		ellipsoid: WGS84(),
	}, nil
}


// Equals checks if the `other` point is equal to this point
//
// Example
// p1 := geod.LatLonEllipsoidal{52.205, 0.119, geod.WGS84()}
// p2 := geod.LatLonEllipsoidal{52.205, 0.119, geod.WGS84()}
// equal := p1.Equals(p2) // true
func (l LatLonEllipsoidal)Equals(other LatLonEllipsoidal) bool {
	epsilon := math.Nextafter(1.0, 2.0)-1.0
        if math.Abs(float64(l.Latitude) - float64(other.Latitude)) > epsilon {
		return false
	}
        if math.Abs(float64(l.Longitude) - float64(other.Longitude)) > epsilon {
		return false
	}
	if math.Abs(l.Height - other.Height) > epsilon {
		return false
	}
	if l.ellipsoid != other.ellipsoid {
		return false
	}
        return true
}
