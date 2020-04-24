package main

// HitRecord defines the intersection of a Ray and an Actor
type HitRecord struct {
	Distance float64
	Position Vec3
	Normal   Vec3
	Material Material
}
