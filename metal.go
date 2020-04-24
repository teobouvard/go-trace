package main

import (
	"github.com/teobouvard/gotrace/util"
)

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
	fuzziness := RandLambertian().Scale(m.fuzz)
	scattered := Ray{record.Position, reflectedDirection.Add(fuzziness)}
	attenuation := m.albedo
	scatters := scattered.Direction.Dot(record.Normal) > 0
	return scatters, attenuation, scattered
}
