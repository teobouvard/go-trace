package main

import (
	"math"
	"math/rand"

	"github.com/teobouvard/gotrace/util"
)

// Material define the way actors interact with a ray
type Material interface {
	Scatter(ray Ray, hit HitRecord) (bool, Vec3, Ray)
}

// Lambertian is a diffuse texture
type Lambertian struct {
	albedo Vec3
}

// Scatter defines how a lambertian material scatters a Ray
func (l Lambertian) Scatter(ray Ray, hit HitRecord) (bool, Vec3, Ray) {
	scatterDirection := hit.Normal.Add(RandSphere())
	scattered := Ray{hit.Position, scatterDirection, ray.Time}
	return true, l.albedo, scattered
}

// Metal is a reflective material
type Metal struct {
	albedo Vec3
	fuzz   float64
}

// NewMetal returns a metal material from its albedo
func NewMetal(albedo Vec3, fuzz float64) Metal {
	return Metal{
		albedo: albedo,
		fuzz:   util.Clamp(fuzz, 0, 1),
	}
}

// Scatter defines the behaviour of rays when they hit Metal material
func (m Metal) Scatter(ray Ray, record HitRecord) (bool, Vec3, Ray) {
	reflectedDirection := ray.Direction.Unit().Reflect(record.Normal)
	fuzziness := RandSphere().Scale(m.fuzz)
	scattered := Ray{record.Position, reflectedDirection.Add(fuzziness), ray.Time}
	attenuation := m.albedo
	scatters := scattered.Direction.Dot(record.Normal) > 0
	return scatters, attenuation, scattered
}

// Dielectric is a glass-like material
type Dielectric struct {
	n float64
}

func shlick(cosine float64, nRatio float64) float64 {
	r0 := math.Pow((1-nRatio)/(1+nRatio), 2)
	return r0 + (1-r0)*math.Pow(1-cosine, 5)

}

// Scatter defines the behaviour of rays when they hit Metal material
func (d Dielectric) Scatter(ray Ray, hit HitRecord) (bool, Vec3, Ray) {
	var (
		outNormal Vec3
		nRatio    float64
		cosTheta  float64
	)
	dot := ray.Direction.Dot(hit.Normal)
	if dot > 0 {
		// ray escapes the material
		outNormal = hit.Normal.Neg()
		nRatio = d.n
		cosTheta = dot / ray.Direction.Norm()
		cosTheta = math.Sqrt(1.0 - d.n*d.n*(1.0-cosTheta*cosTheta))
	} else {
		// ray enters the material
		outNormal = hit.Normal
		nRatio = 1.0 / d.n
		cosTheta = -dot / ray.Direction.Norm()
	}
	incidentDirection := ray.Direction.Unit()
	wasRefracted, refracted := incidentDirection.Refract(outNormal, nRatio)
	var direction Vec3
	if wasRefracted && rand.Float64() >= shlick(cosTheta, nRatio) {
		// refraction possible
		direction = refracted
	} else {
		// reflection
		direction = incidentDirection.Reflect(outNormal)
	}
	return true, WHITE, Ray{hit.Position, direction, ray.Time}
}
