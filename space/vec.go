package space

import (
	"fmt"
	"io"
	"math"
	"math/rand"

	"github.com/teobouvard/gotrace/util"
)

// Vec3 defines a 3-dimensional vector
type Vec3 struct {
	e [3]float64
}

// NewVec3 creates a Vec3 from its components
func NewVec3(i float64, j float64, k float64) Vec3 {
	return Vec3{
		e: [3]float64{i, j, k},
	}
}

// X returns the first component of v
func (v Vec3) X() float64 {
	return v.e[0]
}

// Y returns the first component of v
func (v Vec3) Y() float64 {
	return v.e[1]
}

// Z returns the first component of v
func (v Vec3) Z() float64 {
	return v.e[2]
}

// Add returns the sum of vecs
func Add(vecs ...Vec3) Vec3 {
	var e [3]float64
	for _, v := range vecs {
		e[0] += v.e[0]
		e[1] += v.e[1]
		e[2] += v.e[2]
	}
	return Vec3{
		e: e,
	}
}

// Neg returns the opposite of v
func Neg(v Vec3) Vec3 {
	return Scale(v, -1)
}

// Scale returns v scaled by t
func Scale(v Vec3, t float64) Vec3 {
	return NewVec3(v.e[0]*t, v.e[1]*t, v.e[2]*t)
}

// Div returns the scaling of v by 1/t
func Div(v Vec3, t float64) Vec3 {
	if t == 0 {
		panic("Division by zero")
	}
	return Scale(v, 1/t)
}

// Dot returns the dot (inner) product between u and v
func Dot(u Vec3, v Vec3) float64 {
	return u.e[0]*v.e[0] + u.e[1]*v.e[1] + u.e[2]*v.e[2]
}

// Cross returns the cross product between v1 and v2
func Cross(u Vec3, v Vec3) Vec3 {
	return NewVec3(
		u.e[1]*v.e[2]-u.e[2]*v.e[1],
		u.e[2]*v.e[0]-u.e[0]*v.e[2],
		u.e[0]*v.e[1]-u.e[1]*v.e[0],
	)
}

// Mul returns the *termwise* product between v1 and v2
func Mul(u Vec3, v Vec3) Vec3 {
	return NewVec3(
		u.e[0]*v.e[0],
		u.e[1]*v.e[1],
		u.e[2]*v.e[2],
	)
}

// Reflect computes the reflection of v if it hits a surface of normal n
func Reflect(v Vec3, n Vec3) Vec3 {
	return Add(v, Neg(Scale(n, 2*Dot(v, n))))
}

// Unit returns a unit vector from v
func Unit(v Vec3) Vec3 {
	return Div(v, v.Norm())
}

// Norm returns the euclidean norm of v
func (v Vec3) Norm() float64 {
	return math.Sqrt(v.SquareNorm())
}

// SquareNorm returns the square of the euclidean norm of v
func (v Vec3) SquareNorm() float64 {
	return v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2]
}

// RandUnitSphere returns a random vector inside the unit sphere
func RandUnitSphere() Vec3 {
	r := rand.Float64()
	theta := rand.Float64() * math.Pi
	phi := 2 * rand.Float64() * math.Pi
	x := r * math.Sin(theta) * math.Cos(phi)
	y := r * math.Sin(theta) * math.Sin(phi)
	z := r * math.Cos(theta)
	return NewVec3(x, y, z)
}

// RandLambertian returns vector drawn from a lambertian distribution inside the unit sphere
func RandLambertian() Vec3 {
	a := 2.0 * rand.Float64() * math.Pi
	z := 2.0 * (rand.Float64() - 0.5)
	r := math.Sqrt(1 - z*z)
	return NewVec3(r*math.Cos(a), r*math.Sin(a), z)
}

// WriteColor writes the color of v to f
func (v Vec3) WriteColor(f io.Writer, samples int) {
	colorRange := 256.0
	minv := 0.0
	maxv := 0.999
	scale := 1.0 / float64(samples)
	r := math.Sqrt(scale * v.e[0])
	g := math.Sqrt(scale * v.e[1])
	b := math.Sqrt(scale * v.e[2])

	ir := int(colorRange * util.Clamp(r, minv, maxv))
	ig := int(colorRange * util.Clamp(g, minv, maxv))
	ib := int(colorRange * util.Clamp(b, minv, maxv))
	fmt.Fprintf(f, "%v %v %v\n", ir, ig, ib)
}

func (v Vec3) String() string {
	return fmt.Sprintf("%v %v %v", v.e[0], v.e[1], v.e[2])
}
