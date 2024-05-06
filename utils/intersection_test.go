package utils_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/starboard-nz/go-geodesy/utils"
	"github.com/starboard-nz/orb"
)

func TestSegmentIntersection(t *testing.T) {
	t.Run("Simple Intersection", func(t *testing.T) {
		p1 := orb.Point{0, 0}
		p2 := orb.Point{20, 20}
		q1 := orb.Point{10, 0}
		q2 := orb.Point{10, 20}

		const δ = 0.0001
		is := utils.SegmentIntersection(p1, p2, q1, q2)
		require.NotNil(t, is)
		assert.InDeltaf(t, is[0], 10, δ, "Longitude: %f", is[0])
		assert.InDeltaf(t, is[1], 10.15589, δ, "Latitude: %f", is[0])
	})

	t.Run("Parallel", func(t *testing.T) {
		p1 := orb.Point{0, 0}
		p2 := orb.Point{20, 20}
		q1 := orb.Point{10, 0}
		q2 := orb.Point{30, 20}

		is := utils.SegmentIntersection(p1, p2, q1, q2)
		require.Nil(t, is)
	})

	t.Run("Colinear", func(t *testing.T) {
		p1 := orb.Point{0, 0}
		p2 := orb.Point{20, 20}
		q1 := orb.Point{10, 10}
		q2 := orb.Point{30, 30}

		is := utils.SegmentIntersection(p1, p2, q1, q2)
		require.NotNil(t, is)
	})

	t.Run("Colinear sharing a point", func(t *testing.T) {
		p1 := orb.Point{0, 0}
		p2 := orb.Point{20, 20}
		q1 := orb.Point{20, 20}
		q2 := orb.Point{40, 40}

		is := utils.SegmentIntersection(p1, p2, q1, q2)
		require.NotNil(t, is)
	})

	t.Run("Sharing a point", func(t *testing.T) {
		p1 := orb.Point{0, 0}
		p2 := orb.Point{20, 20}
		q1 := orb.Point{20, 20}
		q2 := orb.Point{40, 60}

		is := utils.SegmentIntersection(p1, p2, q1, q2)
		require.NotNil(t, is)
	})

	t.Run("Middle of one line intersects point of the other", func(t *testing.T) {
		p1 := orb.Point{0, 0}
		p2 := orb.Point{20, 20}
		q1 := orb.Point{20, 0}
		q2 := orb.Point{20, 40}

		is := utils.SegmentIntersection(p1, p2, q1, q2)
		require.NotNil(t, is)
	})

	t.Run("No intersection", func(t *testing.T) {
		p1 := orb.Point{0, 0}
		p2 := orb.Point{20, 20}
		q1 := orb.Point{10, 0}
		q2 := orb.Point{10, 5}

		is := utils.SegmentIntersection(p1, p2, q1, q2)
		require.Nil(t, is)
	})

	t.Run("Intersection on AM crossing lines", func(t *testing.T) {
		p1 := orb.Point{170, 10}
		p2 := orb.Point{-170, -10}
		q1 := orb.Point{-170, 10}
		q2 := orb.Point{170, -10}

		is := utils.SegmentIntersection(p1, p2, q1, q2)
		require.NotNil(t, is)
	})

	t.Run("Intersection on AM crossing lines", func(t *testing.T) {
		p1 := orb.Point{170, 10}
		p2 := orb.Point{-170, -20}
		q1 := orb.Point{-170, 20}
		q2 := orb.Point{170, -10}

		is := utils.SegmentIntersection(p1, p2, q1, q2)
		require.NotNil(t, is)
	})

	t.Run("Intersection on AM crossing lines", func(t *testing.T) {
		p1 := orb.Point{170, 10}
		p2 := orb.Point{-170, -10}
		q1 := orb.Point{-175, 10}
		q2 := orb.Point{-175, -10}

		is := utils.SegmentIntersection(p1, p2, q1, q2)
		require.NotNil(t, is)
	})
}

func TestSegmentsIntersect(t *testing.T) {
	t.Run("Simple Intersection", func(t *testing.T) {
		p1 := orb.Point{0, 0}
		p2 := orb.Point{20, 20}
		q1 := orb.Point{10, 0}
		q2 := orb.Point{10, 20}

		is := utils.SegmentsIntersect(p1, p2, q1, q2)
		require.True(t, is)
	})

	t.Run("Parallel", func(t *testing.T) {
		p1 := orb.Point{0, 0}
		p2 := orb.Point{20, 20}
		q1 := orb.Point{10, 0}
		q2 := orb.Point{30, 20}

		is := utils.SegmentsIntersect(p1, p2, q1, q2)
		require.False(t, is)
	})

	t.Run("Colinear sharing a point", func(t *testing.T) {
		p1 := orb.Point{0, 0}
		p2 := orb.Point{20, 20}
		q1 := orb.Point{20, 20}
		q2 := orb.Point{40, 40}

		is := utils.SegmentsIntersect(p1, p2, q1, q2)
		require.True(t, is)
	})

	t.Run("No intersection", func(t *testing.T) {
		p1 := orb.Point{0, 0}
		p2 := orb.Point{20, 20}
		q1 := orb.Point{10, 0}
		q2 := orb.Point{10, 5}

		is := utils.SegmentsIntersect(p1, p2, q1, q2)
		require.False(t, is)
	})
}

func BenchmarkSegmentIntersection(b *testing.B) {
	const N = 100000
	testP1 := make([]orb.Point, N)
	testP2 := make([]orb.Point, N)
	testP3 := make([]orb.Point, N)
	testP4 := make([]orb.Point, N)
	for i := 0; i < N; i++ {
		testP1[i] = orb.Point{
			rand.Float64() * 100, // nolint:gosec
			rand.Float64() * 100, // nolint:gosec
		}
		testP2[i] = orb.Point{
			rand.Float64() * 100, // nolint:gosec
			rand.Float64() * 100, // nolint:gosec
		}
		testP3[i] = orb.Point{
			rand.Float64() * 100, // nolint:gosec
			rand.Float64() * 100, // nolint:gosec
		}
		testP4[i] = orb.Point{
			rand.Float64() * 100, // nolint:gosec
			rand.Float64() * 100, // nolint:gosec
		}
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = utils.SegmentIntersection(testP1[n%N], testP2[n%N], testP3[n%N], testP4[n%N])
	}
}

func BenchmarkSegmentsIntersect(b *testing.B) {
	const N = 100000
	testP1 := make([]orb.Point, N)
	testP2 := make([]orb.Point, N)
	testP3 := make([]orb.Point, N)
	testP4 := make([]orb.Point, N)
	for i := 0; i < N; i++ {
		testP1[i] = orb.Point{
			rand.Float64() * 100, // nolint:gosec
			rand.Float64() * 100, // nolint:gosec
		}
		testP2[i] = orb.Point{
			rand.Float64() * 100, // nolint:gosec
			rand.Float64() * 100, // nolint:gosec
		}
		testP3[i] = orb.Point{
			rand.Float64() * 100, // nolint:gosec
			rand.Float64() * 100, // nolint:gosec
		}
		testP4[i] = orb.Point{
			rand.Float64() * 100, // nolint:gosec
			rand.Float64() * 100, // nolint:gosec
		}
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = utils.SegmentsIntersect(testP1[n%N], testP2[n%N], testP3[n%N], testP4[n%N])
	}
}
