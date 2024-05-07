package geod_test

/**
 * Copyright (c) 2024, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	geod "github.com/starboard-nz/go-geodesy"
)

func TestPlanar(t *testing.T) {
	δ := 0.0001

	t.Run("InitialBearingTo (not crossing the antimeridian)", func(t *testing.T) {
		p1 := geod.NewLatLonPlanar(20, 20)

		testData := map[geod.LatLon]float64{
			geod.NewLatLon(50, 20):                     0,
			geod.NewLatLon(20, 50):                     90,
			geod.NewLatLon(-10, 20):                    180,
			geod.NewLatLon(20, -10):                    270,
			geod.NewLatLon(50, 50):                     45,
			geod.NewLatLon(-10, 50):                    135,
			geod.NewLatLon(-10, -10):                   225,
			geod.NewLatLon(50, -10):                    315,
			geod.NewLatLon(20, 20):                     math.NaN(),
			geod.NewLatLon(20+math.Sqrt(1600-400), 40): 30,
			geod.NewLatLon(0, 20+math.Sqrt(1600-400)):  120,
			geod.NewLatLon(20-math.Sqrt(1600-400), 0):  210,
			geod.NewLatLon(0, 20-math.Sqrt(1600-400)):  240,
			geod.NewLatLon(40, 20-math.Sqrt(1600-400)): 300,
		}

		for ll, exp := range testData {
			b := p1.InitialBearingTo(ll)
			assert.InDeltaf(t, exp, float64(b), δ, "result: (%f, %f) -> %v", ll.Latitude, ll.Longitude, b)
		}
	})

	t.Run("InitialBearingTo (crossing the antimeridian)", func(t *testing.T) {
		p1 := geod.NewLatLonPlanar(20, 170)

		testData := map[geod.LatLon]float64{
			geod.NewLatLon(50, 170):                      0,
			geod.NewLatLon(20, -160):                     90,
			geod.NewLatLon(-10, 170):                     180,
			geod.NewLatLon(20, 140):                      270,
			geod.NewLatLon(50, -160):                     45,
			geod.NewLatLon(-10, -160):                    135,
			geod.NewLatLon(-10, 140):                     225,
			geod.NewLatLon(50, 140):                      315,
			geod.NewLatLon(20, 170):                      math.NaN(),
			geod.NewLatLon(20+math.Sqrt(1600-400), -170): 30,
			geod.NewLatLon(0, -190+math.Sqrt(1600-400)):  120,
			geod.NewLatLon(20-math.Sqrt(1600-400), 150):  210,
			geod.NewLatLon(0, 170-math.Sqrt(1600-400)):   240,
			geod.NewLatLon(40, 170-math.Sqrt(1600-400)):  300,
		}

		for ll, exp := range testData {
			b := p1.InitialBearingTo(ll)
			assert.InDeltaf(t, exp, float64(b), δ, "result: (%f, %f) -> %v", ll.Latitude, ll.Longitude, b)
		}
	})

	t.Run("InitialBearingTo (crossing the antimeridian) 2", func(t *testing.T) {
		p1 := geod.NewLatLonPlanar(20, -160)

		testData := map[geod.LatLon]float64{
			geod.NewLatLon(50, -160):                     0,
			geod.NewLatLon(20, -130):                     90,
			geod.NewLatLon(-10, -160):                    180,
			geod.NewLatLon(20, 170):                      270,
			geod.NewLatLon(50, -130):                     45,
			geod.NewLatLon(-10, -130):                    135,
			geod.NewLatLon(-10, 170):                     225,
			geod.NewLatLon(50, 170):                      315,
			geod.NewLatLon(20, -160):                     math.NaN(),
			geod.NewLatLon(20+math.Sqrt(1600-400), -140): 30,
			geod.NewLatLon(0, -160+math.Sqrt(1600-400)):  120,
			geod.NewLatLon(20-math.Sqrt(1600-400), 180):  210,
			geod.NewLatLon(20-math.Sqrt(1600-400), -180): 210,
			geod.NewLatLon(0, -160-math.Sqrt(1600-400)):  240,
			geod.NewLatLon(40, -160-math.Sqrt(1600-400)): 300,
		}

		for ll, exp := range testData {
			b := p1.InitialBearingTo(ll)
			assert.InDeltaf(t, exp, float64(b), δ, "result: (%f, %f) -> %v", ll.Latitude, ll.Longitude, b)
		}
	})

	t.Run("Distance (not crossing the antimeridian)", func(t *testing.T) {
		OctagonSth := geod.NewLatLonPlanar(-45.8745, 170.5033)
		OctagonNth := geod.NewLatLon(-45.8736, 170.5038)
		dist := OctagonSth.DistanceTo(OctagonNth)
		assert.InDeltaf(t, float64(107.2692), float64(dist.Metre()), δ, "distance: %v", dist)

		p1 := geod.NewLatLonPlanar(20, 165)
		p2 := geod.NewLatLon(20, 175)
		dist = p1.DistanceTo(p2)
		assert.InDeltaf(t, float64(104489*10), float64(dist.Metre()), δ, "distance: %v", dist)
	})

	t.Run("Distance (crossing the antimeridian)", func(t *testing.T) {
		p1 := geod.NewLatLonPlanar(20, 175)
		p2 := geod.NewLatLon(20, -175)
		dist := p1.DistanceTo(p2)
		assert.InDeltaf(t, float64(104489*10), float64(dist.Metre()), δ, "distance: %v", dist)

		p1 = geod.NewLatLonPlanar(20, -175)
		p2 = geod.NewLatLon(20, 175)
		dist = p1.DistanceTo(p2)
		assert.InDeltaf(t, float64(104489*10), float64(dist.Metre()), δ, "distance: %v", dist)
	})

	t.Run("IntermediatePointTo", func(t *testing.T) {
		type testPoints struct {
			p0 geod.LatLonPlanar
			p1 geod.LatLon
			fr float64
		}

		testData := map[testPoints]geod.LatLon{
			{p0: geod.NewLatLonPlanar(-10, 10), p1: geod.LatLon{Latitude: 30, Longitude: 50}, fr: 0.5}:  {Latitude: 10, Longitude: 30},
			{p0: geod.NewLatLonPlanar(-10, 10), p1: geod.LatLon{Latitude: 30, Longitude: 50}, fr: 0.25}: {Latitude: 0, Longitude: 20},
			{p0: geod.NewLatLonPlanar(0, 170), p1: geod.LatLon{Latitude: 0, Longitude: -150}, fr: 0.50}: {Latitude: 0, Longitude: -170},
			{p0: geod.NewLatLonPlanar(0, 170), p1: geod.LatLon{Latitude: 0, Longitude: -150}, fr: 0.25}: {Latitude: 0, Longitude: 180},
			{p0: geod.NewLatLonPlanar(0, 180), p1: geod.LatLon{Latitude: 0, Longitude: -180}, fr: 0.25}: {Latitude: 0, Longitude: 180},
			{p0: geod.NewLatLonPlanar(0, -150), p1: geod.LatLon{Latitude: 0, Longitude: 170}, fr: 0.50}: {Latitude: 0, Longitude: -170},
			{p0: geod.NewLatLonPlanar(0, -150), p1: geod.LatLon{Latitude: 0, Longitude: 170}, fr: 0.75}: {Latitude: 0, Longitude: -180},
			{p0: geod.NewLatLonPlanar(0, -150), p1: geod.LatLon{Latitude: 0, Longitude: 170}, fr: 0.80}: {Latitude: 0, Longitude: 178},
		}

		for d, exp := range testData {
			mp := d.p0.IntermediatePointTo(d.p1, d.fr)
			assert.Equal(t, exp, mp, "%v -> %v", d.p0, d.p1)
		}
	})
}
