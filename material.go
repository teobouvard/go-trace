package gotrace

import (
	"math"

	"github.com/teobouvard/gotrace/util"
)

/*
Material define the way actors interact with a ray

Scatter

@in

	ray : an incident ray
	hit : the record for the hit of the ray with a geometry

@out

	bool : true if the material scatters the ray
	Vec3 : the attenuation of the scattered ray
	Ray : the scattered ray
*/
type Material interface {
	Scatter(ray Ray, hit HitRecord) (bool, Vec3, Ray)
	Emit(u, v float64, pos Vec3) Vec3
}

// Lambertian is a diffuse material
type Lambertian struct {
	albedo Texture
}

// Scatter defines how a lambertian material scatters a Ray
func (l Lambertian) Scatter(ray Ray, hit HitRecord) (bool, Vec3, Ray) {
	scatterDirection := hit.Normal.Add(RandSphere(ray.RandSource))
	scattered := Ray{hit.Position, scatterDirection, ray.Time, ray.RandSource}
	attenuation := l.albedo.Value(hit.U, hit.V, hit.Position)
	return true, attenuation, scattered
}

// Emit defines how a Lambertian emits light (it doesn't)
func (l Lambertian) Emit(u, v float64, pos Vec3) Vec3 {
	return BLACK
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
	fuzziness := RandSphere(ray.RandSource).Scale(m.fuzz)
	scattered := Ray{record.Position, reflectedDirection.Add(fuzziness), ray.Time, ray.RandSource}
	attenuation := m.albedo
	scatters := scattered.Direction.Dot(record.Normal) > 0
	return scatters, attenuation, scattered
}

// Emit defines how a Metal emits light (it doesn't)
func (m Metal) Emit(u, v float64, pos Vec3) Vec3 {
	return BLACK
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
	if wasRefracted && ray.RandSource.Float64() >= shlick(cosTheta, nRatio) {
		// refraction possible + shlick probability
		direction = refracted
	} else {
		// reflection
		direction = incidentDirection.Reflect(outNormal)
	}
	return true, WHITE, Ray{hit.Position, direction, ray.Time, ray.RandSource}
}

// Emit defines how a lambertian emits light (it doesn't)
func (d Dielectric) Emit(u, v float64, pos Vec3) Vec3 {
	return BLACK
}

// DiffuseLight is a light-emitting material
type DiffuseLight struct {
	emit Texture
}

// Scatter implements the scatter interface for a DiffuseLight material
func (l DiffuseLight) Scatter(ray Ray, hit HitRecord) (bool, Vec3, Ray) {
	return false, Vec3{}, Ray{}
}

// Emit implements the emit interface for a DiffuseLight material
func (l DiffuseLight) Emit(u, v float64, pos Vec3) Vec3 {
	return l.emit.Value(u, v, pos)
}

// Isotropic is a material scattering in random direction
type Isotropic struct {
	albedo Texture
}

// Scatter of isotropic material scatters a new ray in a random direction at the hit
func (i Isotropic) Scatter(ray Ray, hit HitRecord) (bool, Vec3, Ray) {
	return true,
		i.albedo.Value(hit.U, hit.V, hit.Position),
		Ray{hit.Position, RandSphere(ray.RandSource), ray.Time, ray.RandSource}
}

// Emit defines how an isotropic material doesn't emit light
func (i Isotropic) Emit(u float64, v float64, pos Vec3) Vec3 {
	return BLACK
}
