package gotrace

import (
	"image"
	"image/draw"
	"log"
	"math"
	"os"

	"github.com/ojrac/opensimplex-go"
	"github.com/teobouvard/gotrace/util"
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
	freq float64
}

// Value implements the texture interface for a CheckerTexture
func (t CheckerTexture) Value(u, v float64, pos Vec3) Vec3 {
	sines := math.Sin(t.freq*pos.X) * math.Sin(t.freq*pos.Y) * math.Sin(t.freq*pos.Z)
	if sines < 0.0 {
		return t.odd.Value(u, v, pos)
	}
	return t.even.Value(u, v, pos)
}

// Noise is an opensimplex noise
type Noise struct {
	noise     opensimplex.Noise
	frequency float64
}

// Value implements the texture interface for a Noise Texture
func (t Noise) Value(u, v float64, pos Vec3) Vec3 {
	scaled := pos.Scale(t.frequency)
	sample := t.noise.Eval3(scaled.X, scaled.Y, scaled.Z)
	return WHITE.Scale(0.5 * (1.0 + sample))
}

// Marble is a marble-like texture
type Marble struct {
	noise      opensimplex.Noise
	depth      int
	turbulence float64
	scale      float64
}

// genTurbulence creates a turbulence effect by summing noise at different frequencies
func (t Marble) genTurbulence(pos Vec3) float64 {
	sum := 0.0
	freq := pos
	weight := 1.0
	for i := 0; i < t.depth; i++ {
		sum += weight * t.noise.Eval3(t.scale*freq.X, t.scale*freq.Y, t.scale*freq.Z)
		weight *= 0.5
		freq = freq.Scale(2)
	}
	return math.Abs(sum)
}

// Value implements the texture interface for a Marble texture
func (t Marble) Value(u, v float64, pos Vec3) Vec3 {
	turbulence := t.genTurbulence(pos)
	return WHITE.Scale(0.5 * (1.0 + math.Sin(t.scale*pos.Y+t.turbulence*turbulence)))
}

// Image is a texture mapped to an image file
type Image struct {
	data    image.Image
	xoffset float64
	yoffset float64
}

// NewImage creates an image texture from the path to the image, and an offset on the x axis
// The offset is given as a percentage of the width
func NewImage(file string, xoffset, yoffset float64) Image {
	f, err := os.Open(file)
	src, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	bounds := src.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(img, img.Bounds(), src, bounds.Min, draw.Src)
	return Image{img, xoffset / 100.0, yoffset / 100.0}
}

// Value implements the texture interface for an Image texture
func (t Image) Value(u, v float64, pos Vec3) Vec3 {
	width := t.data.Bounds().Max.X - 1
	height := t.data.Bounds().Max.Y - 1
	x := int(util.Map(u, 0, 1, 0, float64(width)))
	y := int(util.Map(v, 0, 1, 0, float64(height)))
	x = int(float64(x)+t.xoffset*float64(width)) % width
	y = int(float64(y)+t.yoffset*float64(height)) % height
	color := t.data.At(x, height-y)
	r, g, b, _ := color.RGBA()
	return Vec3{float64(r) / 65535, float64(g) / 65535, float64(b) / 65535}
}
