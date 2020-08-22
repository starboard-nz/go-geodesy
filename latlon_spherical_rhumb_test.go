package geod

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

import (
	"testing"
	"math"
)

func TestRhumb(t *testing.T) {
	p1 := NewLatLonRhumb(51.127, 1.338)
	p2 := NewLatLon(50.964, 1.853)
	dist := p1.DistanceTo(p2)
	if math.Round(dist.Metres()) != 40308 {
		t.Errorf("Incorrect result")
	}

	brng := p1.InitialBearingTo(p2)
	if math.Round(10 * float64(brng)) != 1167 {
		t.Errorf("Incorrect result")
	}

	dest := p1.DestinationPoint(40300, Degrees(116.7))
	if dest.Latitude.RoundTo(4) != 50.9642 || dest.Longitude.RoundTo(4) != 1.8530 {
		t.Errorf("Incorrect result")
	}

	mp := p1.MidPointTo(p2)
	if mp.Latitude.RoundTo(4) != 51.0455 || mp.Longitude.RoundTo(4) != 1.5957 {
		t.Errorf("Incorrect result")
	}
}
