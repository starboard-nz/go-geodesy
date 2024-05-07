package geod_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	geod "github.com/starboard-nz/go-geodesy"
)

func TestMercator(t *testing.T) {
	const δ = 0.000001

	testData := map[geod.LatLon]geod.MercatorPoint{
		geod.LatLon{Latitude: 0, Longitude: 0}:                      geod.MercatorPoint{Y: 0.5, X: 0.5},
		geod.LatLon{Latitude: 0, Longitude: -180}:                   geod.MercatorPoint{Y: 0.5, X: 0},
		geod.LatLon{Latitude: 0, Longitude: 180}:                    geod.MercatorPoint{Y: 0.5, X: 1},
		geod.LatLon{Latitude: geod.MercatorMaxLat, Longitude: 180}:  geod.MercatorPoint{Y: 1, X: 1},
		geod.LatLon{Latitude: -geod.MercatorMaxLat, Longitude: 180}: geod.MercatorPoint{Y: 0, X: 1},
		geod.LatLon{Latitude: -45, Longitude: 0}:                    geod.MercatorPoint{Y: 0.359725, X: 0.5},
		geod.LatLon{Latitude: 45, Longitude: 0}:                     geod.MercatorPoint{Y: 0.640275, X: 0.5},
	}

	for ll, exp := range testData {
		mp := ll.MercatorPoint()
		assert.InDeltaf(t, exp.X, mp.X, δ, "X coordinate of LatLon{%f, %f} - expected %f got %f",
			ll.Latitude, ll.Longitude, exp.X, mp.X)
		assert.InDeltaf(t, exp.Y, mp.Y, δ, "Y coordinate of LatLon{%f, %f} - expected %f got %f",
			ll.Latitude, ll.Longitude, exp.Y, mp.Y)
	}
}

func TestMercatorConversionAndBack(t *testing.T) {
	const δ = 0.000001
	for i := 0; i < 10; i++ {
		ll := geod.LatLon{
			Longitude: geod.Degrees(rand.Float64()*360 - 180), // nolint:gosec
			Latitude:  geod.Degrees(rand.Float64()*160 - 80),  // nolint:gosec
		}

		mp := ll.MercatorPoint()
		ll1 := mp.LatLon()

		assert.InDelta(t, float64(ll.Longitude), float64(ll1.Longitude), δ)
		assert.InDelta(t, float64(ll.Latitude), float64(ll1.Latitude), δ)
	}
}

func BenchmarkMercator(b *testing.B) {
	const N = 100000
	testPoints := make([]geod.LatLon, N)
	for i := 0; i < N; i++ {
		testPoints[i] = geod.LatLon{
			Longitude: geod.Degrees(rand.Float64()*360 - 180), // nolint:gosec
			Latitude:  geod.Degrees(rand.Float64()*160 - 80),  // nolint:gosec
		}
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = testPoints[n%N].MercatorPoint()
	}
}

func BenchmarkInverseMercator(b *testing.B) {
	const N = 100000
	testPoints := make([]geod.MercatorPoint, N)
	for i := 0; i < N; i++ {
		testPoints[i] = geod.MercatorPoint{
			X: rand.Float64(), // nolint:gosec
			Y: rand.Float64(), // nolint:gosec
		}
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = testPoints[n%N].LatLon()
	}
}
