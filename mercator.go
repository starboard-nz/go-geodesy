package geod

import (
	"math"

	"github.com/starboard-nz/orb"
)

type MercatorPoint struct {
	X float64
	Y float64
}

const π = math.Pi
const MercatorMaxLat = Degrees(85.05112877980644)

// MercatorPoint converts the Latitude/Longitude pair to a X/Y coordinates using Mercator projection.
// The resulting coordinates will be in the [0..1] range, so for rendering images, multiply by the
// horizontal and vertical resolution.
// Latitudes over MercatorMaxLat are not supported.
func (ll LatLon) MercatorPoint() MercatorPoint {
	if ll.Latitude > MercatorMaxLat {
		return MercatorPoint{math.NaN(), math.NaN()}
	}

	x := (float64(ll.Longitude) + 180) / 360

	// convert from degrees to radians
	latRad := ll.Latitude.Radians()

	y := 0.5 + math.Log(math.Tan(π/4 + latRad/2))/(2*π)

	return MercatorPoint{X: x, Y: y}
}

// MercatorPoint convert a point in Mercator projection the a Latitude/Longitude.
// The Mercator coordinates must be in the [0..1] range, so divide by the horizontal/vertical resolution.
func (mp MercatorPoint) LatLon() LatLon {
	latRad := 2*(math.Atan(math.Exp((mp.Y - 0.5) * 2*π))-π/4)
	lat := DegreesFromRadians(latRad)
	lon := Degrees(mp.X * 360 - 180)

	return LatLon{Latitude: lat, Longitude: lon}
}

func MultiPolygonToMercator(mp orb.MultiPolygon) orb.MultiPolygon {
	return nil
}
