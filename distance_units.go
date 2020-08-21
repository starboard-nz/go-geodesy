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

// DistanceUnits represents a distance between 2 points.
// Use Metres() or Kilometres() to get the distance in the unit of your choice.
// If you prefer imperial units, use NauticalMiles() Miles() or Feet().
// Or if you're in the US but like SI standards, you may want to use Meters() or Kilometers()  :)
type DistanceUnits float64

// Valid returns true if the distance is valid. Invalid distances are returned by
// functions when the result cannot be calculated.
func (d DistanceUnits)Valid() bool {
	return !math.IsNaN(float64(d))
}

// Metres returns the DistanceUnits d in metres
func (d DistanceUnits)Metres() float64 {
	return float64(d)
}

// Meters also returns the DistanceUnits d in metres, but in US English
func (d DistanceUnits)Meters() float64 {
	return float64(d)
}

// Kilometres returns the DistanceUnits d in kilometres
func (d DistanceUnits)Kilometres() float64 {
	return float64(d) / 1000.0
}

// Kilometers also returns the DistanceUnits d in kilometres, but in US English
func (d DistanceUnits)Kilometers() float64 {
	return float64(d) / 1000.0
}

// NauticalMiles returns the DistanceUnits d in, you guessed it, nautical miles
func (d DistanceUnits)NauticalMiles() float64 {
	return float64(d) / 1852.0
}

// Miles returns the DistanceUnits d in miles
func (d DistanceUnits)Miles() float64 {
	return float64(d) / 1609.344
}

// Feet returns the DistanceUnits d in feet
func (d DistanceUnits)Feet() float64 {
	return float64(d) / 0.3048
}
