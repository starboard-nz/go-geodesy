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
)

// Vector3D represents a 3 dimensional vector
type Vector3D struct {
	X, Y, Z float64
}

// Length returns the length (magnitude or norm) of the vector.
func (v Vector3D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Copy returns an identical copy of the vector
func (v Vector3D) Copy() Vector3D {
	return Vector3D{v.X, v.Y, v.Z}
}

// Equals returns true if the vector equals the `other` vector, false otherwise
func (v Vector3D) Equals(other Vector3D) bool {
	if v.X == other.X && v.Y == other.Y && v.Z == other.Z {
		return true
	}
	return false
}

// Str returns a string representation of the vector, rounded to 3 decimal points
func (v Vector3D) Str() string {
	return fmt.Sprintf("[%.3f,%.3f,%.3f]", v.X, v.Y, v.Z)
}

// Plus adds the `other` vector to the vector
// Returns a copy of the resulting vector.
func (v Vector3D) Plus(other Vector3D) Vector3D {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	return v
}

// Minus subtracts the `other` vector from the vector
// Returns a copy of the resulting vector.
func (v Vector3D) Minus(other Vector3D) Vector3D {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	return v
}

// Times multiplies the vector by a scalar value
// Returns a copy of the multiplied vector.
func (v Vector3D) Times(f float64) Vector3D {
	v.X *= f
	v.Y *= f
	v.Z *= f
	return v
}

// DividedBy divides the vector by a scalar value
// Returns a copy of the divided vector.
func (v Vector3D) DividedBy(f float64) Vector3D {
	v.X /= f
	v.Y /= f
	v.Z /= f
	return v
}

// Dot multiplies the vector by the `other` vector using dot (scalar) product
func (v Vector3D) Dot(other Vector3D) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross multiplies the vector by the `other` vector using cross (vector) product,
// returns the resulting vector
func (v Vector3D) Cross(other Vector3D) Vector3D {
	x := v.Y*other.Z - v.Z*other.Y
	y := v.Z*other.X - v.X*other.Z
	z := v.X*other.Y - v.Y*other.X
	v.X = x
	v.Y = y
	v.Z = z
	return v
}

// Negate negates a vector to point in the opposite direction,
// returns the resulting vector
func (v Vector3D) Negate() Vector3D {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	return v
}

// Unit normalizes a vector to its unit vector, returns the resulting vector
func (v Vector3D) Unit() Vector3D {
	len := v.Length()
	if len == 1 || len == 0 {
		return v
	}
	v.X /= len
	v.Y /= len
	v.Z /= len
	return v
}

// AngleTo calculates the angle between the vector and the `other` vector atan2(|p₁×p₂|, p₁·p₂) or if
// (extra-planar) `n` is not nil then atan2(n·p₁×p₂, p₁·p₂).
//
// Arguments:
//
// `other` - Vector whose angle is to be determined from the `v` vector
// `n` - Plane normal: if not nil, angle is signed +ve if `v` is clockwise looking along `n`, -ve in opposite direction
//
// Returns the angle (in radians) between the `v` vector and the `other` vector in range 0..π if n is nil,
// or range -π..+π if n  is not nil.
func (v Vector3D) AngleTo(other Vector3D, n *Vector3D) float64 {
	sign := 1.0
	if n != nil {
		if v.Cross(other).Dot(*n) < 0 {
			sign = -1.0
		}
	}
	sinθ := v.Cross(other).Length() * sign
	cosθ := v.Dot(other)

	return math.Atan2(sinθ, cosθ)
}

// RotateAround rotates the vector around an axis by a specified angle
//
// Arguments:
//
// `axis` - The axis being rotated around.
// `angle` - The angle of rotation (in degrees)
//
// Returns the rotated vector
func (v Vector3D) RotateAround(axis Vector3D, angle Degrees) Vector3D {
	θ := angle.Radians()

	// en.wikipedia.org/wiki/Rotation_matrix#Rotation_matrix_from_axis_and_angle
	// en.wikipedia.org/wiki/Quaternions_and_spatial_rotation#Quaternion-derived_rotation_matrix
	p := v.Unit()
	a := axis.Unit()

	s := math.Sin(θ)
	c := math.Cos(θ)
	t := 1 - c
	x := a.X
	y := a.Y
	z := a.Z

	// rotation matrix for rotation about supplied axis
	r := [][]float64{
		{t*x*x + c, t*x*y - s*z, t*x*z + s*y},
		{t*x*y + s*z, t*y*y + c, t*y*z - s*x},
		{t*x*z - s*y, t*y*z + s*x, t*z*z + c},
	}

	// multiply r × p
	rp := []float64{
		r[0][0]*p.X + r[0][1]*p.Y + r[0][2]*p.Z,
		r[1][0]*p.X + r[1][1]*p.Y + r[1][2]*p.Z,
		r[2][0]*p.X + r[2][1]*p.Y + r[2][2]*p.Z,
	}

	return Vector3D{rp[0], rp[1], rp[2]}
	// qv en.wikipedia.org/wiki/Rodrigues'_rotation_formula...
}
