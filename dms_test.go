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

func TestParseDMS(t *testing.T) {
	variations := []string{
		"0.0°",
		"0°",
		"000 00 00 ",
		"000°00′00″",
		"000°00′00.0″",
		"0",
	}
	for _, s := range(variations) {
		dd, err := ParseDMS(s)
		if err != nil {
			t.Errorf("ParseDMS failed: %v", err)
		}
		if dd != 0.0 {
			t.Errorf("Invalid result: expected 0, got %v", dd)
		}

		dd, err = ParseDMS(fmt.Sprintf("%sE", s))
		if err != nil {
			t.Errorf("ParseDMS failed: %v", err)
		}
		if dd != 0.0 {
			t.Errorf("Invalid result: expected 0, got %v", dd)
		}

		dd, err = ParseDMS(fmt.Sprintf(" %sE ", s))
		if err != nil {
			t.Errorf("ParseDMS failed: %v", err)
		}
		if dd != 0.0 {
			t.Errorf("Invalid result: expected 0, got %v", dd)
		}

		dd, err = ParseDMS(fmt.Sprintf("	%s ", s))
		if err != nil {
			t.Errorf("ParseDMS failed: %v", err)
		}
		if dd != 0.0 {
			t.Errorf("Invalid result: expected 0, got %v", dd)
		}

		dd, err = ParseDMS(fmt.Sprintf("%s S", s))
		if err != nil {
			t.Errorf("ParseDMS failed: %v", err)
		}
		if dd != 0.0 {
			t.Errorf("Invalid result: expected 0, got %v", dd)
		}

		dd, err = ParseDMS(fmt.Sprintf("-%s", s))
		if err != nil {
			t.Errorf("ParseDMS failed: %v", err)
		}
		if dd != 0.0 {
			t.Errorf("Invalid result: expected 0, got %v", dd)
		}
	}

        variations = []string{
            `45.76260`,
            `45.76260 `,
            `45.76260°`,
            `45°45.756′`,
            `45° 45.756′`,
            `45 45.756`,
            `45°45′45.36″`,
            `45º45'45.36"`,
            `45°45’45.36”`,
            `45 45 45.36 `,
            `45° 45′ 45.36″`,
            `45º 45' 45.36"`,
            `45° 45’ 45.36”`,
        }

	for _, s := range(variations) {
		dd, err := ParseDMS(s)
		if err != nil {
			t.Errorf("ParseDMS failed: %v", err)
		}
		if dd != 45.76260 {
			t.Errorf("Invalid result: expected 0, got %v", dd)
		}

		dd, err = ParseDMS(fmt.Sprintf("%sE", s))
		if err != nil {
			t.Errorf("ParseDMS failed: %v", err)
		}
		if dd != 45.76260 {
			t.Errorf("Invalid result: expected 0, got %v", dd)
		}

		dd, err = ParseDMS(fmt.Sprintf("%sS", s))
		if err != nil {
			t.Errorf("ParseDMS failed: %v", err)
		}
		if dd != -45.76260 {
			t.Errorf("Invalid result: expected 0, got %v", dd)
		}

		dd, err = ParseDMS(fmt.Sprintf("%sN", s))
		if err != nil {
			t.Errorf("ParseDMS failed: %v", err)
		}
		if dd != 45.76260 {
			t.Errorf("Invalid result: expected 0, got %v", dd)
		}

		dd, err = ParseDMS(fmt.Sprintf("%sW", s))
		if err != nil {
			t.Errorf("ParseDMS failed: %v", err)
		}
		if dd != -45.76260 {
			t.Errorf("Invalid result: expected 0, got %v", dd)
		}

		dd, err = ParseDMS(fmt.Sprintf("-%s", s))
		if err != nil {
			t.Errorf("ParseDMS failed: %v", err)
		}
		if dd != -45.76260 {
			t.Errorf("Invalid result: expected 0, got %v", dd)
		}
	}

	dd, err := ParseDMS("45S")
	if err != nil {
		t.Errorf("ParseDMS failed: %v", err)
	}
	if dd != -45.0 {
		t.Errorf("Invalid result: expected 0, got %v", dd)
	}

	dd, err = ParseDMS("45′N")
	if err != nil {
		t.Errorf("ParseDMS failed: %v", err)
	}
	if dd != 45.0/60.0 {
		t.Errorf("Invalid result: expected 0, got %v", dd)
	}

	invalids := []string {
		"0 0 0 0",
		"xxx",
		"",
		"true",
	}
	for _, s := range(invalids) {
		dd, err = ParseDMS(s)
		if err == nil {
			t.Errorf("Should have failed for %q, got %f", s, dd)
		}
	}
}

func TestToDMS(t *testing.T) {
	s := FormatDMS(0, FormatDeg, -1)
	if s != "000.0000°" {
		t.Errorf("Invalid result")
	}
	s = FormatDMS(0, FormatDegMinSec, -1)
	if s != "000°00′00″" {
		t.Errorf("Invalid result")
	}
	s = FormatDMS(0, FormatDegMinSec, 2)
	if s != "000°00′00.00″" {
		t.Errorf("Invalid result")
	}
	
        s = FormatDMS(9.1525, FormatDeg, -1)
	if s != "009.1525°" {
		t.Errorf("Invalid result")
	}
	s = FormatDMS(9.1525, FormatDegMin, -1)
	if s != "009°09.15′" {
		t.Errorf("Invalid result")
	}
	s = FormatDMS(9.1525, FormatDegMinSec, -1)
	if s != "009°09′09″" {
		t.Errorf("Invalid result")
	}
	s = FormatDMS(9.1525, FormatDeg, 6)
	if s != "009.152500°" {
		t.Errorf("Invalid result")
	}
        s = FormatDMS(9.1525, FormatDegMin, 4)
	if s != "009°09.1500′" {
		t.Errorf("Invalid result")
	}
        s = FormatDMS(9.1525, FormatDegMinSec, 2)
	if s != "009°09′09.00″" {
		t.Errorf("Invalid result")
	}
        s = FormatDMS(9.1525, 999, -1)
	if s != "009.1525°" {
		t.Errorf("Invalid result")
	}
        s = FormatDMS(9.1525, 999, 6)
	if s != "009.152500°" {
		t.Errorf("Invalid result")
	}

	// test rounding
	s = FormatDMS(51.99999999999999, FormatDeg, -1)
	if s != "052.0000°" {
		t.Errorf("Invalid result")
	}
        s = FormatDMS(51.99999999999999, FormatDegMin, -1)
	if s != "052°00.00′" {
		t.Errorf("Invalid result")
	}
        s = FormatDMS(51.99999999999999, FormatDegMinSec, -1)
	if s != "052°00′00″" {
		t.Errorf("Invalid result")
	}
        s = FormatDMS(51.19999999999999, FormatDeg, -1)
	if s != "051.2000°" {
		t.Errorf("Invalid result")
	}
        s = FormatDMS(51.19999999999999, FormatDegMin, -1)
	if s != "051°12.00′" {
		t.Errorf("Invalid result")
	}
        s = FormatDMS(51.19999999999999, FormatDegMinSec, -1)
	if s != "051°12′00″" {
		t.Errorf("Invalid result")
	}
}

func TestWrap360(t *testing.T) {
	testValues := map[float64]float64{
		-450: 270,
		-405: 315,
		-360:   0,
		-315:  45,
		-270:  90,
		-225: 135,
		-180: 180,
		-135: 225,
		-90: 270,
		-45: 315,
		0:   0,
		45:  45,
		90:  90,
		135: 135,
		180: 180,
		225: 225,
		270: 270,
		315: 315,
		360:   0,
		405:  45,
		450:  90,
	}
	for k, v := range(testValues) {
		if float64(Wrap360(Degrees(k))) != v {
			t.Errorf("Invalid result for %v: expected %v got %v", k, v, Wrap360(Degrees(k)))
		}
	}
}

func TestWrap180(t *testing.T) {
	testValues := map[float64]float64{
		-450:  -90,
		-405:  -45,
		-360:    0,
		-315:   45,
		-270:   90,
		-225:  135,
		-180: -180,
		-135: -135,
		-90:  -90,
		-45:  -45,
		0:    0,
		45:   45,
		90:   90,
		135:  135,
		180:  180,
		225: -135,
		270:  -90,
		315:  -45,
		360:    0,
		405:   45,
		450:   90,
	}
	for k, v := range(testValues) {
		if float64(Wrap180(Degrees(k))) != v {
			t.Errorf("Invalid result for %v: expected %v got %v", k, v, Wrap180(Degrees(k)))
		}
	}
}

func TestWrap90(t *testing.T) {
	testValues := map[float64]float64{
		-450:  -90,
		-405:  -45,
		-360:    0,
		// -315: 45 TODO: fix!
		-270:   90,
		-225:   45,
		-180:    0,
		-135:  -45,
		-90:  -90,
		-45:  -45,
		0:    0,
		45:   45,
		90:   90,
		135:   45,
		180:    0,
		225:  -45,
		270:  -90,
		315:  -45,
		360:    0,
		405:   45,
		450:   90,
	}
	for k, v := range(testValues) {
		if float64(Wrap90(Degrees(k))) != v {
			t.Errorf("Invalid result for %v: expected %v got %v", k, v, Wrap90(Degrees(k)))
		}
	}
}
