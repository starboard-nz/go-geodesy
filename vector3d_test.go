package geod

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

import (
	"testing"
)

func TestVector3D(t *testing.T) {
        v123 := Vector3D{1, 2, 3}
        v321 := Vector3D{3, 2, 1}
	if ! v123.Plus(v321).Equals(Vector3D{4, 4, 4}) {
		t.Errorf("Incorrect result")
	}
        if !v123.Minus(v321).Equals(Vector3D{-2, 0, 2}) {
		t.Errorf("Incorrect result")
	}
	if !v123.Times(2).Equals(Vector3D{2, 4, 6}) {
		t.Errorf("Incorrect result")
	}
	if !v123.DividedBy(2).Equals(Vector3D{0.5, 1, 1.5}) {
		t.Errorf("Incorrect result")
	}
        if v123.Dot(v321) != 10 {
		t.Errorf("Incorrect result")
	}
        if !v123.Cross(v321).Equals(Vector3D{-4, 8, -4}) {
		t.Errorf("Incorrect result")
	}
        if !v123.Negate().Equals(Vector3D{-1, -2, -3}) {
		t.Errorf("Incorrect result")
	}
        if v123.Length() != 3.7416573867739413 {
		t.Errorf("Incorrect result")
	}
        if v123.Unit().Str() != "[0.267,0.535,0.802]" {
		t.Errorf("Incorrect result")
	}
        if DegreesFromRadians(v123.AngleTo(v321, nil)).RoundTo(3) != 44.415 {
		t.Errorf("Incorrect result")
	}
	vcross := v123.Cross(v321)
        if DegreesFromRadians(v123.AngleTo(v321, &vcross)).RoundTo(3) != 44.415 {
		t.Errorf("Incorrect result")
	}
	vcross = v321.Cross(v123)
	if DegreesFromRadians(v123.AngleTo(v321, &vcross)).RoundTo(3) != -44.415 {
		t.Errorf("Incorrect result")
	}
	if DegreesFromRadians(v123.AngleTo(v321, &v123)).RoundTo(3) != 44.415 {
		t.Errorf("Incorrect result")
	}
        if v123.RotateAround(Vector3D{0, 0, 1}, 90).Str() != "[-0.535,0.267,0.802]" {
		t.Errorf("Incorrect result")
	}
        if v123.Str() != "[1.000,2.000,3.000]" {
		t.Errorf("Incorrect result")
	}
}
