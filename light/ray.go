package light

import (
	"github.com/teobouvard/gotrace/space"
)

// Ray is a light ray
type Ray struct {
	origin    space.Vec3
	direction space.Vec3
}

// NewRay creates a ray with an origin and a direction
func NewRay(origin *space.Vec3, direction *space.Vec3) *Ray {
	return &Ray{
		origin:    *origin,
		direction: *direction,
	}
}

// Origin returns the origin of the ray
func (r *Ray) Origin() *space.Vec3 {
	return &r.origin
}

// Direction returns the direction of the ray
func (r *Ray) Direction() *space.Vec3 {
	return &r.direction
}

// At is looking from the ray origin at the ray direction times t
func (r *Ray) At(t float64) *space.Vec3 {
	return space.Add(&r.origin, space.Mul(&r.direction, t))
}
