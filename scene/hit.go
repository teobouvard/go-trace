package scene

import (
	"github.com/teobouvard/gotrace/space"
)

// HitRecord defines the intersection of a Ray and an Actor
type HitRecord struct {
	dist float64
	pos  space.Vec3
	norm space.Vec3
	mat  Material
}

// NewHitRecord creates a HitRecord
func NewHitRecord(distance float64, position space.Vec3, normal space.Vec3, material Material) HitRecord {
	return HitRecord{
		dist: distance,
		pos:  position,
		norm: normal,
		mat:  material,
	}
}

// Distance returns the distance to the hit
func (r HitRecord) Distance() float64 {
	return r.dist
}

// Normal returns the normal to the hit
func (r HitRecord) Normal() space.Vec3 {
	return r.norm
}

// Position returns the position of the hit
func (r HitRecord) Position() space.Vec3 {
	return r.pos
}

// Material returns the material of the hit
func (r HitRecord) Material() Material {
	return r.mat
}
