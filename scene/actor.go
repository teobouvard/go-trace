package scene

// Actor is an object on the scene
type Actor struct {
	shape   Geometry
	texture Material
}

// NewActor ??
func NewActor(shape Geometry, texture Material) Actor {
	return Actor{
		shape:   shape,
		texture: texture,
	}
}

// Hit ???
func (a Actor) Hit(ray Ray, tMin float64, tMax float64) (bool, HitRecord) {
	if hit, dist, pos, normal := a.shape.Hit(ray, tMin, tMax); hit {
		return true, NewHitRecord(dist, pos, normal, a.texture)
	}
	return false, HitRecord{}
}
