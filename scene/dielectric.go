package scene

import (
	"math"
	"math/rand"

	"github.com/teobouvard/gotrace/space"
)

// Dielectric is a glass-like material
type Dielectric struct {
	n float64
}

// NewDielectric creates a dielectric material from its reflective index
func NewDielectric(n float64) Dielectric {
	return Dielectric{
		n: n,
	}
}

func shlick(cosine float64, nRatio float64) float64 {
	r0 := math.Pow((1-nRatio)/(1+nRatio), 2)
	return r0 + (1-r0)*math.Pow(1-cosine, 5)

}

// Scatter defines the behaviour of rays when they hit Metal material
func (d Dielectric) Scatter(ray Ray, record HitRecord) (scatters bool, attenuation space.Vec3, scattered Ray) {
	scatters = true
	attenuation = space.WHITE // full transparency
	var nRatio float64
	var outNormal space.Vec3
	if space.Dot(ray.Direction(), record.Normal()) < 0 {
		// ray enters the material
		nRatio = 1 / d.n
		outNormal = record.Normal()
	} else {
		// ray escapes the material, the interface normal is negated because it has to be on the inside
		nRatio = d.n
		outNormal = space.Neg(record.Normal())
	}
	unitDirection := space.Unit(ray.Direction())
	cosTheta := space.Dot(space.Neg(unitDirection), outNormal)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)
	if nRatio*sinTheta > 1.0 || rand.Float64() < shlick(cosTheta, nRatio) {
		// above critical angle, full reflection
		reflected := space.Reflect(unitDirection, outNormal)
		scattered = NewRay(record.Position(), reflected)
	} else {
		// refraction possible
		refracted := space.Refract(unitDirection, outNormal, nRatio)
		scattered = NewRay(record.Position(), refracted)
	}
	return
}
