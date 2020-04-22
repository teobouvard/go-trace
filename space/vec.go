package space

import (
	"fmt"
	"io"
	"math"
)

// Vec3 defines a 3-dimensional vector
type Vec3 struct {
	e          [3]float64
	colorRange float64
}

// NewVec3 create a Vec3 from its 3 components
func NewVec3(e0 float64, e1 float64, e2 float64) *Vec3 {
	return &Vec3{
		e:          [3]float64{e0, e1, e2},
		colorRange: 256,
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
func Add(vecs ...*Vec3) *Vec3 {
	var i float64
	var j float64
	var k float64
	for _, v := range vecs {
		i += v.e[0]
		j += v.e[1]
		k += v.e[2]
	}
	return NewVec3(i, j, k)
}

// Sub returns the difference between u and v
func Sub(vecs ...*Vec3) *Vec3 {
	var i float64
	var j float64
	var k float64
	for _, v := range vecs {
		i -= v.e[0]
		j -= v.e[1]
		k -= v.e[2]
	}
	return NewVec3(i, j, k)
}

// Mul returns v scaled by t
func Mul(v *Vec3, t float64) *Vec3 {
	return NewVec3(v.e[0]*t, v.e[1]*t, v.e[2]*t)
}

// Div returns the scaling of v by 1/t
func Div(v *Vec3, t float64) *Vec3 {
	if t == 0 {
		panic("Division by zero")
	}
	return NewVec3(v.e[0]/t, v.e[1]/t, v.e[2]/t)
}

// Dot returns the dot (inner) product between u and v
func Dot(u *Vec3, v *Vec3) float64 {
	return u.e[0]*v.e[0] + u.e[1]*v.e[1] + u.e[2]*v.e[2]
}

// Cross returns the cross product between v1 and v2
func Cross(u *Vec3, v *Vec3) *Vec3 {
	return NewVec3(
		u.e[1]*v.e[2]-u.e[2]*v.e[1],
		u.e[2]*v.e[0]-u.e[0]*v.e[2],
		u.e[0]*v.e[1]-u.e[1]*v.e[0],
	)
}

// Unit returns a unit vector from v
func Unit(v *Vec3) *Vec3 {
	return Div(v, v.Norm())
}

// Norm returns the euclidean norm of v
func (v *Vec3) Norm() float64 {
	return math.Sqrt(v.SquareNorm())
}

// SquareNorm returns the square of the euclidean norm of v
func (v *Vec3) SquareNorm() float64 {
	return v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2]
}

// WriteColor writes the color of v to f
func (v *Vec3) WriteColor(f io.Writer) {
	r := int(v.colorRange * v.e[0])
	g := int(v.colorRange * v.e[1])
	b := int(v.colorRange * v.e[2])
	fmt.Fprintf(f, "%v %v %v\n", r, g, b)
}

func (v *Vec3) String() string {
	return fmt.Sprintf("%v %v %v", v.e[0], v.e[1], v.e[2])
}
