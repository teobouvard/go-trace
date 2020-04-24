package main

import (
	"math"
)

// Camera represents the camera of the scene
type Camera struct {
	origin     Vec3
	horizontal Vec3
	vertical   Vec3
	corner     Vec3
}

// NewCamera creates a camera
func NewCamera(lookFrom Vec3, lookAt Vec3, up Vec3, verticalFOV float64, aspectRatio float64) Camera {
	theta := (math.Pi * verticalFOV) / 180.0
	height := math.Tan(theta / 2.0)
	width := aspectRatio * height

	w := lookFrom.Sub(lookAt).Unit()
	u := up.Cross(w).Unit()
	v := w.Cross(u)

	horizontal := u.Scale(2 * width)
	vertical := v.Scale(2 * height)

	// o - width*u - height*v - w
	corner := lookFrom.Sub(u.Scale(width)).Sub(v.Scale(height)).Sub(w)

	return Camera{
		origin:     lookFrom,
		horizontal: horizontal,
		vertical:   vertical,
		corner:     corner,
	}
}

// RayTo returns the Ray when the camera looks at (u, v)
func (c Camera) RayTo(u float64, v float64) Ray {
	hOffset := c.horizontal.Scale(u)
	vOffset := c.vertical.Scale(v)
	return Ray{c.origin, c.corner.Add(hOffset).Add(vOffset).Sub(c.origin)}
}
