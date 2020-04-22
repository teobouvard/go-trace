package scene

import (
	"github.com/teobouvard/gotrace/space"
)

// Ray is a light ray
type Ray struct {
	o space.Vec3
	d space.Vec3
}

// NewRay creates a ray emitted at origin along direction
func NewRay(origin space.Vec3, direction space.Vec3) Ray {
	return Ray{
		o: origin,
		d: direction,
	}
}

// Origin returns the origin of the ray
func (r Ray) Origin() space.Vec3 {
	return r.o
}

// Direction returns the direction of the ray
func (r Ray) Direction() space.Vec3 {
	return r.d
}

// At is the point of the ray after it travelled t units of time
func (r Ray) At(t float64) space.Vec3 {
	return space.Add(r.o, space.Scale(r.d, t))
}
