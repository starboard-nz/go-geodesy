package geod

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

import (
	"fmt"
	"testing"
)

func TestModel(t *testing.T) {
	p1 := LatLon{10, 20}
	p2 := LatLon{20, 40}
	mp := MidPoint(p1, p2, SphericalModel)
	fmt.Printf("Midpoint: (spherical) %v\n", mp)
	rmp := MidPoint(p1, p2, RhumbModel)
	fmt.Printf("Midpoint: (rhumb) %v\n", rmp)
	mp2 := MidPoint(p1, rmp, RhumbModel)
	mp3 := MidPoint(rmp, p2, RhumbModel)
	mp4 := MidPoint(mp2, mp3, RhumbModel)
	fmt.Printf("Midpoint: (middle of quarters) %v\n", mp4)
	fmt.Printf("Distance: (spherical) %vkm\n", Distance(p1, p2, SphericalModel).Kilometres())
	fmt.Printf("Distance: (rhumb) %vkm\n", Distance(p1, p2, RhumbModel).Kilometres())
	rp1 := NewLatLonRhumb(51.127, 1.338)
	rp2 := LatLon{50.964, 1.853}
	rpm := rp1.MidPointTo(rp2)
	rpq := rp1.MidPointTo(rpm)
	fmt.Printf("Rhumb quarter point: %v\n", rpq)
	fmt.Printf("Rhumb intermediate point: %v\n", rp1.IntermediatePointTo(rp2, 0.25))
}
