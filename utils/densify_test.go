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
	r2 := utils.DensifyRing(ring, e, geod.SphericalModel)
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
}
