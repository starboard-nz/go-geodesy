package utils

import (
	"github.com/starboard-nz/orb"
	geod "github.com/starboard-nz/go-geodesy"
)

// LineStringIntersections finds the intersections of 2 LineStrings (if exists).
func LineStringIntersections(l1, l2 orb.LineString) []orb.Point {
	if len(l1) < 2 || len(l2) < 2 {
		return nil
	}

	var intersections []orb.Point

	for i := 1; i < len(l1); i++ {
		for j := 1; j < len(l2); j++ {
			is := SegmentIntersection(l1[i-1], l1[i], l2[j-1], l2[j])
			if is != nil {
				intersections = append(intersections, *is)
			}
		}
	}

	return intersections
}

// LineStringsIntersect returns true if the 2 LineStrings intersect.
func LineStringsIntersect(l1, l2 orb.LineString) bool {
	if len(l1) < 2 || len(l2) < 2 {
		return false
	}

	for i := 1; i < len(l1); i++ {
		for j := 1; j < len(l2); j++ {
			if SegmentsIntersect(l1[i-1], l1[i], l2[j-1], l2[j]) {
				return true
			}
		}
	}

	return false
}

// SegmentIntersection returns the intersections of 2 segments (p1, p2) and (q1, q2) (if exists).
func SegmentIntersection(p1, p2, q1, q2 orb.Point) *orb.Point {
	var p *orb.Point
	_ = segmentIntersection(p1, p2, q1, q2, &p)

	return p
}

// SegmentsIntersect returns true if segments (p1, p2) and (q1, q2) intersect.
func SegmentsIntersect(p1, p2, q1, q2 orb.Point) bool {
	return segmentIntersection(p1, p2, q1, q2, nil)
}

func segmentIntersection(p1, p2, q1, q2 orb.Point, is **orb.Point) bool {
	var pMin, pMax, qMin, qMax float64

	if p1[0] < p2[0] {
		pMin, pMax = p1[0], p2[0]
	} else {
		pMin, pMax = p2[0], p1[0]
	}
	if q1[0] < q2[0] {
		qMin, qMax = q1[0], q2[0]
	} else {
		qMin, qMax = q2[0], q1[0]
	}

	if pMax < qMin || qMax < pMin {
		return false
	}

	if p1[1] < p2[1] {
		pMin, pMax = p1[1], p2[1]
	} else {
		pMin, pMax = p2[1], p1[1]
	}
	if q1[1] < q2[1] {
		qMin, qMax = q1[1], q2[1]
	} else {
		qMin, qMax = q2[1], q1[1]
	}

	if pMax < qMin || qMax < pMin {
		return false
	}

	mp1 := geod.LatLon{Latitude: geod.Degrees(p1[1]), Longitude: geod.Degrees(p1[0])}.MercatorPoint()
	mp2 := geod.LatLon{Latitude: geod.Degrees(p2[1]), Longitude: geod.Degrees(p2[0])}.MercatorPoint()
	mq1 := geod.LatLon{Latitude: geod.Degrees(q1[1]), Longitude: geod.Degrees(q1[0])}.MercatorPoint()
	mq2 := geod.LatLon{Latitude: geod.Degrees(q2[1]), Longitude: geod.Degrees(q2[0])}.MercatorPoint()

	s1x := mp2.X - mp1.X
	s1y := mp2.Y - mp1.Y
	s2x := mq2.X - mq1.X
	s2y := mq2.Y - mq1.Y

	s := (-s1y*(mp1.X-mq1.X) + s1x*(mp1.Y-mq1.Y)) / (-s2x*s1y + s1x*s2y)
	if !(s >= 0 && s <= 1) {
		return false
	}

	t := (s2x*(mp1.Y-mq1.Y) - s2y*(mp1.X-mq1.X)) / (-s2x*s1y + s1x*s2y)
	if !(t >= 0 && t <= 1) {
		return false
	}

	if is != nil {
		ll := geod.MercatorPoint{X: mp1.X + (t * s1x), Y: mp1.Y + (t * s1y)}.LatLon()
		*is = &orb.Point{float64(ll.Longitude), float64(ll.Latitude)}
	}

	return true
}
