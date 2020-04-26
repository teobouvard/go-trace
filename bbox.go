package gotrace

import "math"

// Bbox is a bounding box of a geometry
type Bbox struct {
	Min Vec3
	Max Vec3
}

// Hit computes the if the intersection of a ray with a bounding box exists
func (b Bbox) Hit(ray Ray, tMin float64, tMax float64) bool {
	origin := ray.Origin.AsArray()
	direction := ray.Direction.AsArray()
	min := b.Min.AsArray()
	max := b.Max.AsArray()
	for i := 0; i < 3; i++ {
		inv := 1.0 / direction[i]
		t0 := (min[i] - origin[i]) * inv
		t1 := (max[i] - origin[i]) * inv
		if inv < 0.0 {
			t0, t1 = t1, t0
		}
		if t0 > tMin {
			tMin = t0
		}
		if t1 < tMax {
			tMax = t1
		}
		// are tMin and tMax reset in each loop ?
		if tMax <= tMin {
			return false
		}
	}
	return true
}

// Merge returns the union of two bounding boxes
func (b Bbox) Merge(o Bbox) Bbox {
	small := Vec3{
		X: math.Min(b.Min.X, o.Min.X),
		Y: math.Min(b.Min.Y, o.Min.Y),
		Z: math.Min(b.Min.Z, o.Min.Z),
	}
	big := Vec3{
		X: math.Max(b.Max.X, o.Max.X),
		Y: math.Max(b.Max.Y, o.Max.Y),
		Z: math.Max(b.Max.Z, o.Max.Z),
	}
	return Bbox{small, big}
}
