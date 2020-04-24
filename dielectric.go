package main

import (
	"math"
	"math/rand"
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
	return true, WHITE, Ray{hit.Position, direction}
}
