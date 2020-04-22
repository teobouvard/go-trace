package scene

import (
	"github.com/teobouvard/gotrace/space"
)

// Lambertian material
type Lambertian struct {
	albedo space.Vec3
}

// NewLambertian creates a lambertian material from its albedo
func NewLambertian(albedo space.Vec3) Lambertian {
	return Lambertian{
		albedo: albedo,
	}
}

// Scatter defines how a lambertian material scatters a Ray
func (l Lambertian) Scatter(ray Ray, record HitRecord) (scatters bool, attenuation space.Vec3, scattered Ray) {
	return
}
