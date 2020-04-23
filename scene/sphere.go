package scene

import (
	"math"

	"github.com/teobouvard/gotrace/space"
)

// Sphere is a spherical actor
type Sphere struct {
	center space.Vec3
	radius float64
}

// NewSphere creates a sphere from its center and its radius
func NewSphere(center space.Vec3, radius float64) Sphere {
	return Sphere{
		center: center,
		radius: radius,
	}
}

// Hit implements the intersection checking of a Ray and the sphere
func (s Sphere) Hit(ray Ray, tMin float64, tMax float64) (hit bool, t float64, pos space.Vec3, normal space.Vec3) {
	oc := space.Add(ray.Origin(), space.Neg(s.center))
	a := ray.Direction().SquareNorm()
	b := space.Dot(oc, ray.Direction())
	c := oc.SquareNorm() - s.radius*s.radius
	discriminant := b*b - a*c
	hit = discriminant > 0

	if hit {
		root := math.Sqrt(discriminant)
		// first solution, closest to camera
		t = (-b - root) / a
		if t < tMax && t > tMin {
			pos = ray.At(t)
			normal = space.Unit(space.Add(pos, space.Neg(s.center)))
			return
		}
		// second solution, farthest from camera
		t = (-b + root) / a
		if t < tMax && t > tMin {
			pos = ray.At(t)
			normal = space.Unit(space.Add(pos, space.Neg(s.center)))
			return
		}
		hit = false
	}
	return
}
