package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/teobouvard/gotrace/util"
)

// Vec3 defines a 3-dimensional vector
type Vec3 struct {
	X, Y, Z float64
}

// Colors
var (
	BLACK = Vec3{0, 0, 0}
	WHITE = Vec3{1, 1, 1}
	RED   = Vec3{1, 0, 0}
	GREEN = Vec3{0, 1, 0}
	BLUE  = Vec3{0, 0, 1}
)

// Add returns u + v
func (u Vec3) Add(v Vec3) Vec3 {
	return Vec3{
		X: u.X + v.X,
		Y: u.Y + v.Y,
		Z: u.Z + v.Z,
	}
}

// Sub returns u-v
func (u Vec3) Sub(v Vec3) Vec3 {
	return Vec3{
		X: u.X - v.X,
		Y: u.Y - v.Y,
		Z: u.Z - v.Z,
	}
}

// Neg returns -v
func (u Vec3) Neg() Vec3 {
	return u.Scale(-1.0)
}

// Scale returns v scaled by t
func (u Vec3) Scale(t float64) Vec3 {
	return Vec3{
		X: t * u.X,
		Y: t * u.Y,
		Z: t * u.Z,
	}
}

// Div returns the scaling of v by 1/t
func (u Vec3) Div(t float64) Vec3 {
	if t == 0 {
		panic("division by zero")
	}
	return Vec3{
		X: u.X / t,
		Y: u.Y / t,
		Z: u.Z / t,
	}
}

// Dot returns the dot (inner) product between u and v
func (u Vec3) Dot(v Vec3) float64 {
	return u.X*v.X + u.Y*v.Y + u.Z*v.Z
}

// Cross returns the cross product between u and v
func (u Vec3) Cross(v Vec3) Vec3 {
	return Vec3{
		X: u.Y*v.Z - u.Z*v.Y,
		Y: u.Z*v.X - u.X*v.Z,
		Z: u.X*v.Y - u.Y*v.X,
	}
}

// Mul returns the *termwise* product betweenu and v
func (u Vec3) Mul(v Vec3) Vec3 {
	return Vec3{
		X: u.X * v.X,
		Y: u.Y * v.Y,
		Z: u.Z * v.Z,
	}

}

// Reflect computes the reflection of v if it hits a surface of normal n
func (u Vec3) Reflect(n Vec3) Vec3 {
	proj := n.Scale(2 * u.Dot(n))
	return u.Sub(proj)
}

// Refract returns the refraction of u at an interface
func (u Vec3) Refract(n Vec3, nRatio float64) (bool, Vec3) {
	dot := u.Dot(n)
	discriminant := 1.0 - nRatio*nRatio*(1-dot*dot)
	if discriminant > 0 {
		refracted := u.Sub(n.Scale(dot)).Scale(nRatio).Sub(n.Scale(math.Sqrt(discriminant)))
		return true, refracted
	}
	return false, Vec3{}
}

// Unit returns a unit vector from u
func (u Vec3) Unit() Vec3 {
	return u.Div(u.Norm())
}

// Norm returns the euclidean norm of u
func (u Vec3) Norm() float64 {
	return math.Sqrt(u.SquareNorm())
}

// SquareNorm returns the square of the euclidean norm of u
func (u Vec3) SquareNorm() float64 {
	return u.Dot(u)
}

// RandSphere returns vector drawn from a lambertian distribution inside the unit sphere
func RandSphere() Vec3 {
	a := 2.0 * rand.Float64() * math.Pi
	z := 2.0 * (rand.Float64() - 0.5)
	r := math.Sqrt(1 - z*z)
	return Vec3{
		r * math.Cos(a),
		r * math.Sin(a),
		z,
	}
}

// RandDisk returns a random vector in the unit disk
func RandDisk() Vec3 {
	theta := 2 * math.Pi * rand.Float64()
	r := rand.Float64()
	return Vec3{X: r * math.Cos(theta), Y: r * math.Sin(theta)}
}

// RandVec returns a random vector with coordinates in [0, 1)
func RandVec() Vec3 {
	return Vec3{rand.Float64(), rand.Float64(), rand.Float64()}
}

// RandVecInterval returns a random vector with coordinates in [low, high)
func RandVecInterval(low float64, high float64) Vec3 {
	x := rand.Float64()*(high-low) + low
	y := rand.Float64()*(high-low) + low
	z := rand.Float64()*(high-low) + low
	return Vec3{x, y, z}
}

// AsArray returns the coordinates of the vector as an array of size 3
func (u Vec3) AsArray() [3]float64 {
	return [3]float64{u.X, u.Y, u.Z}
}

// GetColor retruns the RGBA color of the vector
func (u Vec3) GetColor(samples int) color.RGBA {
	maxColor := 255.0
	scale := 1.0 / float64(samples)
	r := math.Sqrt(scale * u.X) // sqrt for alpha correction
	g := math.Sqrt(scale * u.Y)
	b := math.Sqrt(scale * u.Z)

	ir := uint8(util.Map(r, 0, 1, 0, maxColor))
	ig := uint8(util.Map(g, 0, 1, 0, maxColor))
	ib := uint8(util.Map(b, 0, 1, 0, maxColor))

	return color.RGBA{ir, ig, ib, 255}
}
