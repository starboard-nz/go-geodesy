package geod

// Pure Go re-implementation of https://github.com/chrisveness/geodesy

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

// Ellipsoid parameters
// The only ellipsoid defined is WGS84, for use in utm/mgrs, vincenty, nvector.
type Ellipsoid struct {
	a, b, f float64
}

var wgs84 = Ellipsoid{a: 6378137, b: 6356752.314245, f: 1/298.257223563}

// WGS84 is a standard ellipsoid used in cartography, geodesy, and satellite navigation including GPS
func WGS84() Ellipsoid {
	return wgs84
}
