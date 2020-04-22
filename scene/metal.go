package scene

import (
	"github.com/teobouvard/gotrace/space"
	"github.com/teobouvard/gotrace/util"
)

// Metal is a reflective material
type Metal struct {
	albedo space.Vec3
	fuzz   float64
}

// NewMetal returns a metal material from its albedo
func NewMetal(albedo space.Vec3, fuzz float64) Metal {
	return Metal{
		albedo: albedo,
		fuzz:   util.Clamp(fuzz, 0, 1),
	}
}

// Scatter defines the behaviour of rays when they hit Metal material
func (m Metal) Scatter(ray Ray, record HitRecord) (scatters bool, attenuation space.Vec3, scattered Ray) {
	reflectedDirection := space.Reflect(space.Unit(ray.Direction()), record.Normal())
	fuzziness := space.Scale(space.RandLambertian(), m.fuzz)
	scattered = NewRay(record.Position(), space.Add(reflectedDirection, fuzziness))
	attenuation = m.albedo
	scatters = space.Dot(scattered.Direction(), record.Normal()) > 0
	return
}
