package geod

// Pure Go re-implementation of https://github.com/chrisveness/geodesy

/**
 * Copyright (c) 2020, Xerra Earth Observation Institute
 * All rights reserved. Use is subject to License terms.
 * See LICENSE in the root directory of this source tree.
 */

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// FormatDeg, FormatDegMin and FormatDegMinSec are constants the control how FormatDMS should format the degree value.
const (
	FormatDeg       = iota // degrees
	FormatDegMin           // degrees+minutes
	FormatDegMinSec        // degrees+minutes+seconds
)

const dmsSeparator = 0x202f // U+202F = 'narrow no-break space'
var dmsRE *regexp.Regexp = regexp.MustCompile(
	`^-?(?:([0-9.,]+)(?:[°º]|\s|[nwseNWSE]?$))?\s*(?:([0-9.,]+)(?:[′’']|\s|[nwseNWSE]?$))?\s*(?:([0-9.,]+)[″”"]?)?\s*[nwseNWSE]?$`)

// ParseDMS parses a string representing Degrees-Minutes-Seconds into decimal degrees
// This is very flexible on formats, allowing signed decimal degrees, or deg-min-sec optionally
// suffixed by compass direction (NSEW); a variety of separators are accepted. Examples -3.62,
// '3 37 12W', '3°37′12″W'.
// Example:
// lat := geod.ParseDMS("51° 28′ 40.37″ N")
// lon := geod.ParseDMS("000° 00′ 05.29″ W")
// ll := geod.LatLng{Latitude: lat, Longitude: lng}    <---  51.4779°N, 000.0015°W
func ParseDMS(dms string) (Degrees, error) {
	var err error
	var nanDegrees = Degrees(math.NaN())

	if dms == "" {
		return 0.0, fmt.Errorf("ParseDMS: Empty string")
	}

	// check for signed decimal degrees without NSEW, if so return it directly
	if fl, err := strconv.ParseFloat(dms, 64); err == nil {
		return Degrees(fl), nil
	}

	dms = strings.TrimSpace(dms)
	// strip off any sign or compass dir'n & split out separate d/m/s

	dmsParts := dmsRE.FindStringSubmatch(dms)
	if len(dmsParts) == 0 {
		return nanDegrees, fmt.Errorf("Failed to parse DMS string %q", dms)
	}

	var deg float64
	if dmsParts[1] != "" {
		deg, err = strconv.ParseFloat(dmsParts[1], 64)
		if err != nil {
			return nanDegrees, fmt.Errorf("Failed to parse degrees (%v) in DMS string %q", dmsParts[1], dms)
		}
	}

	var min float64
	if dmsParts[2] != "" {
		min, err = strconv.ParseFloat(dmsParts[2], 64)
		if err != nil {
			return nanDegrees, fmt.Errorf("Failed to parse minutes (%v) in DMS string %q", dmsParts[2], dms)
		}
	}

	var sec float64
	if dmsParts[3] != "" {
		sec, err = strconv.ParseFloat(dmsParts[3], 64)
		if err != nil {
			return nanDegrees, fmt.Errorf("Failed to parse seconds (%v) in DMS string %q", dmsParts[3], dms)
		}
	}

	// and convert to decimal degrees...
	deg += min/60.0 + sec/3600.0

	if strings.HasPrefix(dms, "-") || strings.HasSuffix(dms, "W") || strings.HasSuffix(dms, "S") {
		deg = -deg
	}

	return Degrees(deg), nil
}

// FormatDMS converts decimal degrees to a string in deg/min/sec format
// Degree, prime, double-prime symbols are added, but sign is discarded, though no compass direction is added.
// Degrees are zero-padded to 3 digits; for degrees latitude, use slice [1:] to remove a leading zero.
//
// Arguments:
//
// `deg` - degrees to be formatted as specified.
// `format` - one of FormatDeg, FormatDegMin or FormatDegMinSec (degrees, degrees+minutes, degrees+minutes+seconds)
// `dp` - number of decimal places to use - use -1 for defaults: 4 for d, 2 for dm, 0 for dms.
func FormatDMS(deg Degrees, format, dp int) string {
	degf := float64(deg)
	if math.IsNaN(degf) || math.IsInf(degf, 0) {
		// give up here if we can't make a number from degf
		return ""
	}

	// default values
	if dp == -1 {
		switch format {
		case FormatDeg:
			dp = 4
		case FormatDegMin:
			dp = 2
		case FormatDegMinSec:
			dp = 0
		default:
			format = FormatDeg
			dp = 4
		}
	}

	degf = math.Abs(degf) // unsigned result ready for appending compass dir'n

	var dms string
	switch format {
	case FormatDegMin:
		d := math.Floor(degf)                                                    // get component deg
		m := math.Round(math.Pow10(dp)*math.Mod(degf*60, 60.0)) / math.Pow10(dp) // get component min
		if m == 60.0 {                                                           // check for rounding up
			d++
			m = 0.0
		}
		dpad := 0
		if d < 10 {
			dpad = 2
		} else if d < 100 {
			dpad = 1
		}
		mpad := 0
		if m < 10 {
			mpad = 1
		}
		dms = fmt.Sprintf("%s%s°%s%s′",
			"00"[0:dpad], // left-pad with leading zeros
			strconv.FormatFloat(d, 'f', 0, 64),
			"0"[0:mpad],                         // left-pad with leading zeros (note may include decimals)
			strconv.FormatFloat(m, 'f', dp, 64)) // round/right-pad minutes
	case FormatDegMinSec:
		d := math.Floor(degf)                                                      // get component deg
		m := math.Mod(math.Floor(degf*3600/60), 60.0)                              // get component min
		s := math.Round(math.Pow10(dp)*math.Mod(degf*3600, 60.0)) / math.Pow10(dp) // get component sec
		if s == 60.0 {                                                             // check for rounding up
			m++
			s = 0.0
		}
		if m == 60.0 { // check for rounding up
			d++
			m = 0.0
		}
		dpad := 0
		if d < 10 {
			dpad = 2
		} else if d < 100 {
			dpad = 1
		}
		mpad := 0
		if m < 10 {
			mpad = 1
		}
		spad := 0
		if s < 10 {
			spad = 1
		}
		dms = fmt.Sprintf("%s%s°%s%s′%s%s″",
			"00"[0:dpad], // left-pad with leading zeros
			strconv.FormatFloat(d, 'f', 0, 64),
			"0"[0:mpad], // left-pad with leading zeros
			strconv.FormatFloat(m, 'f', 0, 64),
			"0"[0:spad],                         // left-pad with leading zeros (note may include decimals)
			strconv.FormatFloat(s, 'f', dp, 64)) // round/right-pad minutes
	default: // FormatDeg falls under this as well
		dpad := 0
		if degf < 10 {
			dpad = 2
		} else if degf < 100 {
			dpad = 1
		}
		dms = fmt.Sprintf("%s%s°",
			"00"[0:dpad],                           // left-pad with leading zeros (note may include decimals)
			strconv.FormatFloat(degf, 'f', dp, 64)) // round/right-pad degrees
	}

	return dms
}

// Wrap360 contrains `degrees` to range 0..360 (e.g. for bearings); -1 --> 359, 361 --> 1.
func Wrap360(degrees Degrees) Degrees {
	if 0.0 <= float64(degrees) && float64(degrees) < 360.0 {
		// avoid rounding due to arithmetic ops if within range
		return degrees
	}
	return Degrees(math.Mod(math.Mod(float64(degrees), 360)+360, 360)) // sawtooth wave p:360, a:360
}

// Wrap180 constrains `degrees` to range -180..+180 (e.g. for longitude); -181 --> 179, 181 --> -179.
func Wrap180(degrees Degrees) Degrees {
	if -180.0 < float64(degrees) && float64(degrees) <= 180.0 {
		// avoid rounding due to arithmetic ops if within range
		return degrees
	}
	return Degrees(
		math.Mod(
			float64(degrees)+180.0+360*(math.Floor(math.Abs(float64(degrees)/360.0))+1),
			360.0) - 180.0) // sawtooth wave p:180, a:±180
}

// Wrap90 constrains `degrees` to range -90..+90 (e.g. for latitude); -91 --> -89, 91 --> 89.
func Wrap90(degrees Degrees) Degrees {
	if -90.0 <= float64(degrees) && float64(degrees) <= 90.0 {
		// avoid rounding due to arithmetic ops if within range
		return degrees
	}
	// triangle wave p:360 a:±90 TODO: fix e.g. -315°
	return Degrees(math.Abs(math.Mod(math.Mod(float64(degrees), 360.0)+270.0, 360.0)-180.0) - 90.0)
}
