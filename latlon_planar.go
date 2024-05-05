package geod

/**
 * Copyright (c) 2024, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

import (
	"math"

	"github.com/starboard-nz/units"
)

// LatLonPlanar represents a point used for calculations on a 2-dimensional plane
// Longitudes still go -180 to 180 and wrap around and Longitudes go -90 to 90.
// Works across the antimeridian.
type LatLonPlanar struct {
	ll LatLon
}

// PlanarModel returns a `Model` that wraps geodesy calculations using Planar model (2-dimensional plane)
// Only suitable for short distances.
func PlanarModel(ll LatLon, modelArgs ...interface{}) Model {
	if len(modelArgs) != 0 {
		panic("Invalid number of arguments in call to PlanarModel()")
	}
	return LatLonPlanar{ll: ll}
}

// LatLon converts LatLonPlanar to LatLon
func (lls LatLonPlanar)LatLon() LatLon {
	return lls.ll
}

// NewLatLonPlanar creates a new LatLonPlanar struct
func NewLatLonPlanar(latitude, longitude float64) LatLonPlanar {
	return LatLonPlanar{
		ll: LatLon{
			Latitude: Wrap90(Degrees(latitude)),
			Longitude: Wrap180(Degrees(longitude)),
		},
	}
}

// ParseLatLonPlanar parses a latitude/longitude point from a variety of formats
// See ParseLatLon for details.
func ParseLatLonPlanar(args ...interface{}) (LatLonPlanar, error) {
	ll, err := ParseLatLon(args)
	if err != nil {
		return LatLonPlanar{}, err
	}
	return LatLonPlanar{ll: ll}, nil
}

// approximate distances of 1 degree longitude at each latitude in metres
var lngDistances = map[int]float64{
	90: 20.0, 89: 1941, 88: 3881, 87: 5819,	86: 7756,
	85: 9691, 84: 11623, 83: 13551, 82: 15475, 81: 17395,
	80: 19309, 79: 21217, 78: 23118, 77: 25013, 76: 26900,
	75: 28779, 74: 30649, 73: 32510, 72: 34361, 71: 36201,
	70: 38030, 69: 39848, 68: 41654, 67: 43447, 66: 45227,
	65: 46993, 64: 48744, 63: 50481, 62: 52202, 61: 53908,
	60: 55597, 59: 57269, 58: 58924, 57: 60561, 56: 62179,
	55: 63778, 54: 65358, 53: 66918, 52: 68458, 51: 69977,
	50: 71474, 49: 72950, 48: 74403, 47: 75834, 46: 77242,
	45: 78626, 44: 79986, 43: 81322, 42: 82633, 41: 83919,
	40: 85180, 39: 86414, 38: 87622, 37: 88804, 36: 89958,
	35: 91085, 34: 92184, 33: 93256, 32: 94298, 31: 95312,
	30: 96297, 29: 97253, 28: 98179, 27: 99075, 26: 99941,
	25: 100777, 24: 101581, 23: 102355, 22: 103098, 21: 103809,
	20: 104489, 19: 105137, 18: 105753, 17: 106336, 16: 106887,
	15: 107406, 14: 107892, 13: 108345, 12: 108765, 11: 109152,
	10: 109506, 9: 109826, 8: 110113, 7: 110366, 6: 110586,
	5: 110772, 4: 110924, 3: 111043, 2: 111127, 1: 111178, 0: 111195,
}

func (lls LatLonPlanar) DistanceTo(dest LatLon) units.Distance {
	y0 := float64(Wrap90(lls.ll.Latitude))
	y1 := float64(Wrap90(dest.Latitude))
	dy := math.Abs(y0 - y1) * 111195 // metres

	avgLat := int(math.Round(math.Abs(y0 + y1)/2))
	lngDist, ok := lngDistances[avgLat]
	if !ok {
		lngDist = 111195
	}

	x0 := float64(Wrap180(lls.ll.Longitude))
	x1 := float64(Wrap180(dest.Longitude))
	dx := x0 - x1

	// antimeridian issues
	if math.Abs(float64(dx)) >= 180 {
		if dx < 0 {
			dx += 360
		} else {
			dx -= 360
		}
	}

	dx *= lngDist

	return units.Metre(math.Sqrt(dx*dx + dy*dy))
}

// Returns the initial bearing in Degrees from North (0°..360°)
func (lls LatLonPlanar) InitialBearingTo(ll LatLon) Degrees {
	dx := Wrap180(ll.Longitude) - Wrap180(lls.ll.Longitude)

	// antimeridian issues
	if math.Abs(float64(dx)) >= 180 {
		if dx < 0 {
			dx += 360
		} else {
			dx -= 360
		}
	}

	dy := Wrap90(ll.Latitude) - Wrap90(lls.ll.Latitude)
	if dx == 0 {
		if dy == 0 {
			return Degrees(math.NaN())
		}

		if dy > 0 {
			return Degrees(0)
		}

		return Degrees(180)
	}

	rad := math.Atan(float64(dy)/float64(dx))
	deg := 90 - DegreesFromRadians(rad)

	if dx < 0 {
		deg += 180
	}

	return deg
}

func (lls LatLonPlanar) FinalBearingOn(ll LatLon) Degrees {
	return lls.InitialBearingTo(ll)
}

func (lls LatLonPlanar) DestinationPoint(distance float64, bearing Degrees) LatLon {
	panic("not implemented")
}

func (lls LatLonPlanar) MidPointTo(ll LatLon) LatLon {
	return lls.IntermediatePointTo(ll, 0.5)
}

func (lls LatLonPlanar) IntermediatePointTo(ll LatLon, fraction float64) LatLon {
	dx := Wrap180(ll.Longitude) - Wrap180(lls.ll.Longitude)

	// antimeridian issues
	if math.Abs(float64(dx)) >= 180 {
		if dx < 0 {
			dx += 360
		} else {
			dx -= 360
		}
	}

	dy := Wrap90(ll.Latitude) - Wrap90(lls.ll.Latitude)

	// Planar or not, longitudes don't make sense at the poles
	if ll.Latitude == 90 || ll.Latitude == -90 {
		return LatLon{
			Latitude: Wrap90(lls.ll.Latitude + dy * Degrees(fraction)),
			Longitude: lls.ll.Longitude,
		}
	}

	if lls.ll.Latitude == 90 || lls.ll.Latitude == -90 {
		return LatLon{
			Latitude: Wrap90(lls.ll.Latitude + dy * Degrees(fraction)),
			Longitude: ll.Longitude,
		}
	}

	return LatLon{
		Latitude: Wrap90(lls.ll.Latitude + dy * Degrees(fraction)),
		Longitude: Wrap180(lls.ll.Longitude + dx * Degrees(fraction)),
	}
}

func (lls LatLonPlanar) IntermediatePointsTo(ll LatLon, fractions []float64) []LatLon {
	res := make([]LatLon, 0, len(fractions))

	for _, fr := range fractions {
		res = append(res, lls.IntermediatePointTo(ll, fr))
	}

	return res
}

