package gotrace

import (
	"math"

	"github.com/ojrac/opensimplex-go"
)

/*
Texture interface

Value

@in

u, v : coordinates of the point

@out

Vec3 : color at the given coordinates
*/
type Texture interface {
	Value(u, v float64, pos Vec3) Vec3
}

// ConstantTexture is a uniform texture with a single color
type ConstantTexture struct {
	color Vec3
}

// Value implements the texture interface for a ConstantTexture
func (t ConstantTexture) Value(u, v float64, pos Vec3) Vec3 {
	return t.color
}

// CheckerTexture is a checkboard-like texture
type CheckerTexture struct {
	odd  Texture
	even Texture
}

// Value implements the texture interface for a CheckerTexture
func (t CheckerTexture) Value(u, v float64, pos Vec3) Vec3 {
	freq := 10.0
	sines := math.Sin(freq*pos.X) * math.Sin(freq*pos.Y) * math.Sin(freq*pos.Z)
	if sines < 0.0 {
		return t.odd.Value(u, v, pos)
	}
	return t.even.Value(u, v, pos)
}

// Noise is an opensimplex noise
type Noise struct {
	noise opensimplex.Noise
}

// Value implements the texture interface for a Noise Texture
func (t Noise) Value(u,v float64, pos Vec3) Vec3 {
	return WHITE.Scale(t.noise.Eval3(pos.X, pos.Y, pos.Z))
}