package utils_test

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	geod "github.com/starboard-nz/go-geodesy"
	"github.com/starboard-nz/go-geodesy/utils"
	"github.com/starboard-nz/orb"
	"github.com/starboard-nz/units"
	// "github.com/starboard-nz/orb/geojson"
)

func TestSegmentError(t *testing.T) {
	p0 := orb.Point{-154.5000, -35}
	p1 := orb.Point{-180.0000, -35}
	p2 := orb.Point{-165, -25}

	e := utils.SegmentError(p0, p1, geod.SphericalModel)
	assert.InDelta(t, float64(e.Km()), 75.0483, 0.0001)

	e = utils.SegmentError(p0, p1, geod.RhumbModel)
	assert.Equal(t, 0.0, float64(e.Km()))

	e = utils.SegmentError(p1, p2, geod.RhumbModel)
	assert.InDelta(t, float64(e.Km()), 18.2367, 0.0001)
}

func TestDensifyRing(t *testing.T) {
	p0 := orb.Point{-154.5000, -35}
	p1 := orb.Point{-180.0000, -35}
	p2 := orb.Point{-165, -25}
	ring := orb.Ring{p0, p1, p2, p0}

	e := utils.SegmentError(p0, p1, geod.SphericalModel)
	r2 := utils.DensifyRing(ring, e.Metre(), geod.SphericalModel)
	assert.Equal(t, ring, r2)
	r3 := utils.DensifyRing(ring, units.Metre(float64(e.Metre())-1), geod.SphericalModel)
	assert.Len(t, r3, 5)

	denseRing := utils.DensifyRing(ring, units.Metre(10), geod.RhumbModel)
	assert.Len(t, denseRing, 130)

	/*
		fc := geojson.NewFeatureCollection()
		fc.Append(geojson.NewFeature(denseRing))
		fc.Append(geojson.NewFeature(ring))
		rawJSON, _ := fc.MarshalJSON()
		fmt.Printf("%s\n", rawJSON)
	*/

	// should be the same
	denseAgainRing := utils.DensifyRing(denseRing, units.Metre(10), geod.RhumbModel)
	assert.Equal(t, denseRing, denseAgainRing)

	// Densify AM crossing ring (near equator)
	p0 = orb.Point{170, -10}
	p1 = orb.Point{-170, -10}
	p2 = orb.Point{-170, 10}
	p3 := orb.Point{170, 10}
	ring = orb.Ring{p0, p1, p2, p3, p0}

	e = utils.SegmentError(p0, p1, geod.SphericalModel)
	r2 = utils.DensifyRing(ring, e.Metre(), geod.SphericalModel)
	assert.Equal(t, ring, r2) // Rings are the error as the tolerance value

	r3 = utils.DensifyRing(ring, units.Metre(float64(e.Metre())-1), geod.SphericalModel)
	assert.Len(t, r3, 7) // Ring densified once for both E/W segments using error - 1

	denseRing = utils.DensifyRing(ring, units.Metre(10), geod.RhumbModel)
	assert.Len(t, denseRing, 5)

	// Densify AM crossing ring (not near equator)
	p0 = orb.Point{170, -70}
	p1 = orb.Point{-170, -70}
	p2 = orb.Point{-170, -60}
	p3 = orb.Point{170, -60}
	ring = orb.Ring{p0, p1, p2, p3, p0}

	smallE := utils.SegmentError(p2, p3, geod.SphericalModel)
	bigE := utils.SegmentError(p0, p1, geod.SphericalModel)

	r2 = utils.DensifyRing(ring, smallE.Metre(), geod.SphericalModel)
	assert.Equal(t, ring, r2) // Rings are the same if using the small error as the tolerance value

	r3 = utils.DensifyRing(ring, units.Metre(float64(bigE.Metre())-1), geod.SphericalModel)
	assert.Len(t, r3, 7) // Ring densified once for both E/W segments using big error - 1

	denseRing = utils.DensifyRing(ring, units.Metre(10), geod.RhumbModel)
	assert.Len(t, denseRing, 5)

	// Densify AM crossing ring (not near equator)
	p0 = orb.Point{179, -70}
	p1 = orb.Point{-1, -70}
	p2 = orb.Point{-1, -60}
	p3 = orb.Point{179, -60}
	ring = orb.Ring{p0, p1, p2, p3, p0}

	smallE = utils.SegmentError(p2, p3, geod.SphericalModel)
	bigE = utils.SegmentError(p0, p1, geod.SphericalModel)

	r2 = utils.DensifyRing(ring, smallE.Metre(), geod.SphericalModel)
	assert.Equal(t, ring, r2) // Rings are the same if using the small error as the tolerance value

	r3 = utils.DensifyRing(ring, units.Metre(float64(bigE.Metre())-1), geod.SphericalModel)
	assert.Len(t, r3, 7) // Ring densified once for both E/W segments using big error - 1

	denseRing = utils.DensifyRing(ring, units.Metre(10), geod.RhumbModel)
	assert.Len(t, denseRing, 5) // No densification required for rhumb model

	denseRing = utils.DensifyRing(ring, units.Metre(10), geod.SphericalModel)
	assert.Len(t, denseRing, 76) // Densified appropriately

	// Densify AM crossing ring (not near equator)
	p0 = orb.Point{179, -70}
	p1 = orb.Point{-1, -70}
	p2 = orb.Point{-1, -60}
	p3 = orb.Point{179, -60}
	ring = orb.Ring{p0, p1, p2, p3, p0}

	smallE = utils.SegmentError(p2, p3, geod.SphericalModel)
	bigE = utils.SegmentError(p0, p1, geod.SphericalModel)

	r2 = utils.DensifyRing(ring, smallE.Metre(), geod.SphericalModel)
	assert.Equal(t, ring, r2)

	r3 = utils.DensifyRing(ring, units.Metre(float64(bigE.Metre())-1), geod.SphericalModel)
	assert.Len(t, r3, 7)

	denseRing = utils.DensifyRing(ring, units.Metre(10), geod.RhumbModel)
	assert.Len(t, denseRing, 5) // No densification required for rhumb model

	denseRing = utils.DensifyRing(ring, units.Metre(10), geod.SphericalModel)
	assert.Len(t, denseRing, 76) // Densified appropriately

	// Densify AM crossing ring (-180 to 180)
	p0 = orb.Point{180, -70}
	p1 = orb.Point{0, -70}
	p2 = orb.Point{-180, -70}
	p3 = orb.Point{-180, -60}
	p4 := orb.Point{0, -60}
	p5 := orb.Point{180, -60}
	ring = orb.Ring{p0, p1, p2, p3, p4, p5, p0}

	denseRing = utils.DensifyRing(ring, units.Metre(1), geod.RhumbModel) // No densification required for rhumb model
	assert.Len(t, denseRing, 7)

	denseRing = utils.DensifyRing(ring, units.Metre(1), geod.SphericalModel) // No densification required for rhumb model
	assert.Len(t, denseRing, 175)

	// Densify -180 to 180 rectangle with no intermediate points
	// Note if we calculate the segment error of {180, -10} to {-180, -10} we will get value of 0 which will cause
	// a stack overflow if we pass it to utils.DensifyRing as the tolerance.
	p0 = orb.Point{180, -10}
	p1 = orb.Point{-180, -10}
	p2 = orb.Point{-180, 70}
	p3 = orb.Point{180, 0}
	ring = orb.Ring{p0, p1, p2, p3, p0}

	denseRing = utils.DensifyRing(ring, units.Metre(1), geod.RhumbModel) // No densification required for rhumb model
	assert.Len(t, denseRing, 5)

	denseRing = utils.DensifyRing(ring, units.Metre(1), geod.RhumbModel) // No densification required for spherical model
	assert.Len(t, denseRing, 5)
}
