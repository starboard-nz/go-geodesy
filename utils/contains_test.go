package utils_test

import (
	"testing"

	geod "github.com/starboard-nz/go-geodesy"
	"github.com/starboard-nz/go-geodesy/utils"
	"github.com/starboard-nz/orb"
)

func TestRingContains(t *testing.T) {
	ring := orb.Ring{
		{0, 0}, {0, 1}, {1, 1}, {1, 0.5}, {2, 0.5},
		{2, 1}, {3, 1}, {3, 0}, {0, 0},
	}

	// +-+ +-+
	// | | | |
	// | +-+ |
	// |     |
	// +-----+

	cases := []struct {
		name   string
		point  orb.Point
		result bool
	}{
		{
			name:   "in base",
			point:  orb.Point{1.5, 0.25},
			result: true,
		},
		{
			name:   "in right tower",
			point:  orb.Point{0.5, 0.75},
			result: true,
		},
		{
			name:   "in middle",
			point:  orb.Point{1.5, 0.75},
			result: false,
		},
		{
			name:   "in left tower",
			point:  orb.Point{2.5, 0.75},
			result: true,
		},
		{
			name:   "in tp middle",
			point:  orb.Point{1.5, 1.0},
			result: false,
		},
		{
			name:   "above",
			point:  orb.Point{2.5, 1.75},
			result: false,
		},
		{
			name:   "below",
			point:  orb.Point{2.5, -1.75},
			result: false,
		},
		{
			name:   "left",
			point:  orb.Point{-2.5, -0.75},
			result: false,
		},
		{
			name:   "right",
			point:  orb.Point{3.5, 0.75},
			result: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ring.Reverse()
			val := utils.RingContains(ring, tc.point, false, geod.RhumbModel)

			if val != tc.result {
				t.Errorf("wrong containment: %v != %v", val, tc.result)
			}

			// should not care about orientation
			ring.Reverse()
			val = utils.RingContains(ring, tc.point, false, geod.RhumbModel)
			if val != tc.result {
				t.Errorf("wrong containment: %v != %v", val, tc.result)
			}
		})
	}

	// points should all be in
	for i, p := range ring {
		if !utils.RingContains(ring, p, false, geod.RhumbModel) {
			t.Errorf("point index %d: should be inside", i)
		}
	}

	// on all the segments should be in.
	for i := 1; i < len(ring); i++ {
		c := interpolate(ring[i], ring[i-1], 0.5)
		if !utils.RingContains(ring, c, false, geod.RhumbModel) {
			t.Errorf("index %d centroid: should be inside", i)
		}
	}

	// colinear with segments but outside
	for i := 1; i < len(ring); i++ {
		p := interpolate(ring[i], ring[i-1], 5)
		if utils.RingContains(ring, p, false, geod.RhumbModel) {
			t.Errorf("index %d centroid: should not be inside", i)
		}

		p = interpolate(ring[i], ring[i-1], -5)
		if utils.RingContains(ring, p, false, geod.RhumbModel) {
			t.Errorf("index %d centroid: should not be inside", i)
		}
	}
}

func TestPolygonContains(t *testing.T) {
	// should exclude holes
	p := orb.Polygon{
		{{0, 0}, {3, 0}, {3, 3}, {0, 3}, {0, 0}},
	}

	if !utils.PolygonContains(p, orb.Point{1.5, 1.5}, geod.RhumbModel) {
		t.Errorf("should contain point")
	}

	// ring oriented same as outer ring
	p = append(p, orb.Ring{{1, 1}, {2, 1}, {2, 2}, {1, 2}, {1, 1}})
	if utils.PolygonContains(p, orb.Point{1.5, 1.5}, geod.RhumbModel) {
		t.Errorf("should not contain point in hole")
	}

	p[1].Reverse() // oriented correctly as opposite of outer
	if utils.PolygonContains(p, orb.Point{1.5, 1.5}, geod.RhumbModel) {
		t.Errorf("should not contain point in hole")
	}

	// point is a vertex of the hole
	if !utils.PolygonContains(p, orb.Point{2, 2}, geod.RhumbModel) {
		t.Errorf("should contain point which touches vertex of hole")
	}

	// point touches edge of the hole
	if !utils.PolygonContains(p, orb.Point{2, 1.5}, geod.RhumbModel) {
		t.Errorf("should contain point which touches edge of hole")
	}
}

func TestMultiPolygonContains(t *testing.T) {
	// should exclude holes
	mp := orb.MultiPolygon{
		{{{0, 0}, {1, 0}, {1, 1}, {0, 1}, {0, 0}}},
	}

	if !utils.MultiPolygonContains(mp, orb.Point{0.5, 0.5}, geod.RhumbModel) {
		t.Errorf("should contain point")
	}

	if utils.MultiPolygonContains(mp, orb.Point{1.5, 1.5}, geod.RhumbModel) {
		t.Errorf("should not contain point")
	}

	mp = append(mp, orb.Polygon{{{2, 0}, {3, 0}, {3, 1}, {2, 1}, {2, 0}}})

	if !utils.MultiPolygonContains(mp, orb.Point{2.5, 0.5}, geod.RhumbModel) {
		t.Errorf("should contain point")
	}

	if utils.MultiPolygonContains(mp, orb.Point{1.5, 0.5}, geod.RhumbModel) {
		t.Errorf("should not contain point")
	}

	// Meridian tests
	mp = append(mp, orb.Polygon{{{-10, -10}, {10, -10}, {10, 10}, {-10, 10}, {-10, -10}}})

	if !utils.MultiPolygonContains(mp, orb.Point{10, 10}, geod.RhumbModel) {
		t.Errorf("should contain point")
	}

	if utils.MultiPolygonContains(mp, orb.Point{10.00000000001, 10}, geod.RhumbModel) {
		t.Errorf("should not contain point")
	}

	if !utils.MultiPolygonContains(mp, orb.Point{-9.99999999999999, 10}, geod.RhumbModel) {
		t.Errorf("should contain point")
	}

}

func interpolate(a, b orb.Point, percent float64) orb.Point {
	return orb.Point{
		a[0] + percent*(b[0]-a[0]),
		a[1] + percent*(b[1]-a[1]),
	}
}
