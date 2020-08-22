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

func TestSpherical(t *testing.T) {
	p1 := NewLatLonSpherical(52.205, 0.119)
	p2 := NewLatLon(48.857, 2.351)
	if math.Round(p1.DistanceTo(p2).Metres()) != 404279 {
		t.Errorf("Incorrect result")
	}
	
	SetEarthRadius(3959.0)
	if math.Round(10 * p1.DistanceTo(p2).Metres()) != 2512 {
		t.Errorf("Incorrect result")
	}
	SetEarthRadius(6371000.0)

	brng := p1.InitialBearingTo(p2)
	if math.Round(10 * float64(brng)) != 1562 {
		t.Errorf("Incorrect result")
	}

	brng = p1.FinalBearingOn(p2)
	if math.Round(10 * float64(brng)) != 1579 {
		t.Errorf("Incorrect result")
	}

	mp := p1.MidPointTo(p2)
	if mp.Latitude.RoundTo(4) != 50.5363 || mp.Longitude.RoundTo(4) != 1.2746 {
		t.Errorf("Incorrect result")
	}

	intp := p1.IntermediatePointTo(p2, 0.25)
	if intp.Latitude.RoundTo(4) != 51.3721 || intp.Longitude.RoundTo(4) != 0.7073 {
		t.Errorf("Incorrect result")
	}

	p3 := NewLatLonSpherical(51.47788, -0.00147)
	dp := p3.DestinationPoint(7794, Degrees(300.7))
	if dp.Latitude.RoundTo(4) != 51.5136 || dp.Longitude.RoundTo(4) != -0.0983 {
		t.Errorf("Incorrect result")
	}

	p4 := NewLatLonSpherical(51.8853, 0.2545)
	p5 := NewLatLon(49.0034, 2.5735)
	x := p4.Intersection(108.547, p5, 32.435)
	if x.Latitude.RoundTo(4) != 50.9078 || x.Longitude.RoundTo(4) != 4.5084 {
		t.Errorf("Incorrect result")
	}
}

