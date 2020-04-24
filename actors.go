package main

// HitRecord defines the intersection of a Ray and an Actor
type HitRecord struct {
	Distance float64
	Position Vec3
	Normal   Vec3
	Material Material
}

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

// Collection represents a collection of Actors
type Collection struct {
	actors []Actor
}

// Clear removes all actors from the collection
func (c *Collection) Clear() {
	c.actors = nil
}

// Add appends actors to the collection
func (c *Collection) Add(actors ...Actor) {
	c.actors = append(c.actors, actors...)
}

// Hit returns the closest hit record if an intersection was found
func (c Collection) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	hitAnything := false
	closestHit := tMax
	var closestRecord *HitRecord

	for _, actor := range c.actors {
		// using the closest hit so far as tmax,
		if hit, record := actor.Hit(ray, tMin, closestHit); hit {
			closestRecord = record
			closestHit = record.Distance
			hitAnything = true
		}
	}

	return hitAnything, closestRecord
}
