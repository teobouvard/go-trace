package main

// Actor is an object on the scene having a shape and a texture
type Actor struct {
	shape   Geometry
	texture Material
}

// Hit checks if the geometry is hit by the ray, and creates a hitrecord with its texture
func (a Actor) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	if hit, dist, pos, n := a.shape.Hit(ray, tMin, tMax); hit {
		return true, &HitRecord{
			Distance: dist,
			Position: pos,
			Normal:   n,
			Material: a.texture,
		}
	}
	return false, nil
}
