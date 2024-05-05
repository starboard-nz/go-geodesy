package utils_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	geod "github.com/starboard-nz/go-geodesy"
	"github.com/starboard-nz/go-geodesy/utils"
	"github.com/starboard-nz/orb"
	"github.com/starboard-nz/units"
	"github.com/starboard-nz/orb/geojson"
)

func saveToGeoJSON(fname string, geoms []orb.Geometry, points []orb.Point) error {
	colours := []string{"#ff3300", "#00ff33", "#0033ff", "#6633ee", "#ee6633", "#33ee66"}
	style := func(colour string) geojson.Properties {
		return geojson.Properties {
			"style": map[string]interface{}{
				"color": colour,
				"opacity": 0.7,
				"dashArray": "",
				"weight": 3,
			},
		}
	}

	fc := geojson.NewFeatureCollection()
	for i := range geoms {
		fc.Append(&geojson.Feature{
			Type: "Feature",
			Geometry: geoms[i],
			Properties: style(colours[i%len(colours)]),
		})
	}

	if len(points) > 0 {
		st := style(colours[len(geoms)%len(colours)])
		for i := range points {
			fc.Append(&geojson.Feature{
				Type: "Feature",
				Geometry: points[i],
				Properties: st,
			})
		}
	}

	rawJSON, err := fc.MarshalJSON()
	if err != nil {
		return err
	}

	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s\n", rawJSON)
	if err != nil {
		return err
	}

	return nil
}	

func TestSegmentError(t *testing.T) {
	p0 := orb.Point{-154.5000, -35}
	p1 := orb.Point{-180.0000, -35}
	p2 := orb.Point{-165, -25}

	e := utils.SegmentError(p0, p1, geod.SphericalModel, geod.PlanarModel)
	assert.InDelta(t, float64(e.Km()), 75.0483, 0.0001)

	e = utils.SegmentError(p0, p1, geod.RhumbModel, geod.PlanarModel)
	assert.Equal(t, 0.0, float64(e.Km()))

	e = utils.SegmentError(p1, p2, geod.RhumbModel, geod.PlanarModel)
	assert.InDelta(t, float64(e.Km()), 18.2367, 0.0001)
}

func TestDensifyRing(t *testing.T) {
	t.Run("Simple Spherical", func(t *testing.T) {
		p0 := orb.Point{-154.5000, -35}
		p1 := orb.Point{-180.0000, -35}
		p2 := orb.Point{-165, -25}
		ring := orb.Ring{p0, p1, p2, p0}

		e := utils.SegmentError(p0, p1, geod.SphericalModel, geod.PlanarModel)
		r2, err := utils.DensifyRing(ring, geod.SphericalModel, geod.PlanarModel, e.Metre())
		require.NoError(t, err)
		assert.Equal(t, ring, r2)
		r3, err := utils.DensifyRing(ring, geod.SphericalModel, geod.PlanarModel, e.Metre()-1)
		require.NoError(t, err)
		assert.Len(t, r3, 5)

		denseRing, err := utils.DensifyRing(ring, geod.RhumbModel, geod.PlanarModel, units.Metre(10))
		assert.NoError(t, err)
		assert.Len(t, denseRing, 130, "got %d", len(denseRing))

		t.Logf("Densify Again")

		// should be the same
		denseAgainRing, err := utils.DensifyRing(denseRing, geod.RhumbModel, geod.PlanarModel, units.Metre(10))
		require.NoError(t, err)
		assert.Equal(t, denseRing, denseAgainRing)
		saveToGeoJSON("/tmp/simple-spherical.json", []orb.Geometry{ring, r3}, r3)
		saveToGeoJSON("/tmp/simple-rhumb.json", []orb.Geometry{ring, denseRing}, denseRing)
	})

	t.Run("Densify AM crossing ring (near equator)", func(t *testing.T) {
		p0 := orb.Point{170, -10}
		p1 := orb.Point{-170, -10}
		p2 := orb.Point{-170, 10}
		p3 := orb.Point{170, 10}
		ring := orb.Ring{p0, p1, p2, p3, p0}

		e := utils.SegmentError(p0, p1, geod.SphericalModel, geod.PlanarModel)
		r2, err := utils.DensifyRing(ring, geod.SphericalModel, geod.PlanarModel, e.Metre())
		require.NoError(t, err)
		assert.Equal(t, ring, r2) // Rings are the same if using the small error as the tolerance value

		r3, err := utils.DensifyRing(ring, geod.SphericalModel, geod.PlanarModel, e.Metre()-1)
		require.NoError(t, err)
		assert.Len(t, r3, 7) // Ring densified once for both E/W segments using error - 1

		denseRing, err := utils.DensifyRing(ring, geod.RhumbModel, geod.PlanarModel, units.Metre(10))
		require.NoError(t, err)
		assert.Len(t, denseRing, 5)
		saveToGeoJSON("/tmp/am-crossing1-spherical.json", []orb.Geometry{ring, r3}, r3)
		saveToGeoJSON("/tmp/am-crossing1-rhumb.json", []orb.Geometry{ring, denseRing}, denseRing)
	})

	t.Run("Densify AM crossing ring (not near equator)", func(t *testing.T) {
		p0 := orb.Point{170, -70}
		p1 := orb.Point{-170, -70}
		p2 := orb.Point{-170, -60}
		p3 := orb.Point{170, -60}
		ring := orb.Ring{p0, p1, p2, p3, p0}

		smallE := utils.SegmentError(p2, p3, geod.SphericalModel, geod.PlanarModel)
		bigE := utils.SegmentError(p0, p1, geod.SphericalModel, geod.PlanarModel)

		r2, err := utils.DensifyRing(ring, geod.SphericalModel, geod.PlanarModel, smallE.Metre())
		require.NoError(t, err)
		assert.Equal(t, ring, r2) // Rings are the same if using the small error as the tolerance value

		r3, err := utils.DensifyRing(ring, geod.SphericalModel, geod.PlanarModel, bigE.Metre()-1)
		require.NoError(t, err)
		assert.Len(t, r3, 7) // Ring densified once for both E/W segments using big error - 1

		denseRing, err := utils.DensifyRing(ring, geod.RhumbModel, geod.PlanarModel, units.Metre(10))
		require.NoError(t, err)
		assert.Len(t, denseRing, 5)
		saveToGeoJSON("/tmp/am-crossing2-spherical.json", []orb.Geometry{ring, r3}, r3)
		saveToGeoJSON("/tmp/am-crossing2-rhumb.json", []orb.Geometry{ring, denseRing}, denseRing)
	})

	t.Run("Densify AM crossing ring (not near equator)", func(t *testing.T) {
		p0 := orb.Point{179, -70}
		p1 := orb.Point{-1, -70}
		p2 := orb.Point{-1, -60}
		p3 := orb.Point{179, -60}
		ring := orb.Ring{p0, p1, p2, p3, p0}

		smallE := utils.SegmentError(p2, p3, geod.SphericalModel, geod.PlanarModel)
		bigE := utils.SegmentError(p0, p1, geod.SphericalModel, geod.PlanarModel)

		r2, err := utils.DensifyRing(ring, geod.SphericalModel, geod.PlanarModel, smallE.Metre())
		require.NoError(t, err)
		assert.Equal(t, ring, r2) // Rings are the same if using the small error as the tolerance value

		r3, err := utils.DensifyRing(ring, geod.SphericalModel, geod.PlanarModel, bigE.Metre()-1)
		require.NoError(t, err)
		assert.Len(t, r3, 7) // Ring densified once for both E/W segments using big error - 1

		denseRing, err := utils.DensifyRing(ring, geod.RhumbModel, geod.PlanarModel, units.Metre(10))
		require.NoError(t, err)
		assert.Len(t, denseRing, 5) // No densification required for rhumb model

		denseRing2, err := utils.DensifyRing(ring, geod.RhumbModel, geod.SphericalModel, units.Metre(10))
		assert.NoError(t, err)
		assert.Len(t, denseRing2, 1539) // Densified appropriately
		saveToGeoJSON("/tmp/am-crossing3-spherical1.json", []orb.Geometry{ring, r3}, r3)
		saveToGeoJSON("/tmp/am-crossing3-spherical2.json", []orb.Geometry{ring, denseRing2}, denseRing2)
		saveToGeoJSON("/tmp/am-crossing3-rhumb.json", []orb.Geometry{ring, denseRing}, denseRing)
	})

	t.Run("Densify AM crossing ring (-180 to 180)", func(t *testing.T) {
		p0 := orb.Point{180, -70}
		p1 := orb.Point{0, -70}
		p2 := orb.Point{-180, -70}
		p3 := orb.Point{-180, -60}
		p4 := orb.Point{0, -60}
		p5 := orb.Point{180, -60}
		ring := orb.Ring{p0, p1, p2, p3, p4, p5, p0}

		// No densification required for rhumb model
		denseRing, err := utils.DensifyRing(ring, geod.RhumbModel, geod.PlanarModel, units.Metre(1))
		require.NoError(t, err)
		assert.Len(t, denseRing, 7)

		// No densification required for rhumb model
		denseRing, err = utils.DensifyRing(ring, geod.RhumbModel, geod.SphericalModel, units.Metre(1000))
		require.NoError(t, err)
		assert.Len(t, denseRing, 259)
		saveToGeoJSON("/tmp/am-crossing4-spherical.json", []orb.Geometry{ring, denseRing}, denseRing)
	})
	
	
	t.Run("-180 to 180 rectangle with no intermediate points", func(t *testing.T) {
		p0 := orb.Point{180, -10}
		p1 := orb.Point{-180, -10}
		p2 := orb.Point{-180, 70}
		p3 := orb.Point{180, 0}
		ring := orb.Ring{p0, p1, p2, p3, p0}

		// No densification required for rhumb model
		denseRing, err := utils.DensifyRing(ring, geod.RhumbModel, geod.PlanarModel, units.Metre(1))
		require.NoError(t, err)
		assert.Len(t, denseRing, 5)

		// No densification required for spherical model
		denseRing, err = utils.DensifyRing(ring, geod.RhumbModel, geod.PlanarModel, units.Metre(1))
		require.NoError(t, err)
		assert.Len(t, denseRing, 5)
	})
}

func TestDensifyErrors(t *testing.T) {
	p0 := orb.Point{-154.5000, -55}
	p1 := orb.Point{-180.0000, -35}
	p2 := orb.Point{-165, -25}
	ring := orb.Ring{p0, p1, p2, p0}

	_, err := utils.DensifyRing(ring, geod.SphericalModel, geod.PlanarModel, units.Metre(0))
	assert.Error(t, err, utils.ErrInvalidTolerance)

	denseRing, err := utils.DensifyRing(ring, geod.SphericalModel, geod.PlanarModel, units.Metre(0.0001))
	assert.ErrorIs(t, err, utils.ErrToleranceTooLow)
	assert.Len(t, denseRing, 49153)

	denseRing, err = utils.DensifyRing(ring, geod.SphericalModel, geod.PlanarModel, units.Metre(0.01))
	assert.NoError(t, err)
	assert.Len(t, denseRing, 14499, "Got %v", len(denseRing))
}
