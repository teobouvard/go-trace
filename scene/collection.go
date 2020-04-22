package scene

// Collection represents a collection of Actors
type Collection struct {
	actors []Actor
}

// NewCollection builds a collection from actors
func NewCollection(actors ...Actor) Collection {
	return Collection{
		actors: actors,
	}
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
func (c Collection) Hit(ray Ray, tMin float64, tMax float64) (bool, HitRecord) {
	hitAnything := false
	closestHit := tMax
	var closestRecord HitRecord

	for _, actor := range c.actors {
		if hit, record := actor.Hit(ray, tMin, closestHit); hit {
			closestRecord = record
			closestHit = record.Distance()
			hitAnything = true
		}
	}

	return hitAnything, closestRecord
}
