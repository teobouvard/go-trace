package gotrace

import (
	"math/rand"
	"sort"
)

// HitRecord defines the intersection of a Ray and an Actor
type HitRecord struct {
	Distance float64
	Position Vec3
	Normal   Vec3
	Material Material
	U, V     float64
}

// Actor is an object on the scene having a shape and a material
type Actor struct {
	shape    Geometry
	material Material
}

// Hit checks if the geometry is hit by the ray, and creates a hitrecord with its material
func (a Actor) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	if hit, record := a.shape.Hit(ray, tMin, tMax); hit {
		record.Material = a.material
		return true, record
	}
	return false, nil
}

// Bound returns the bounding box of an actor
func (a Actor) Bound(startTime float64, endTime float64) (bool, *Bbox) {
	return a.shape.Bound(startTime, endTime)
}

// Collection represents a collection of Actors
type Collection []Actor

// Add appends actors to the collection
func (c *Collection) Add(actors ...Actor) {
	*c = append(*c, actors...)
}

// Hit returns the closest hit record if an intersection was found
func (c Collection) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	hitAnything := false
	closestHit := tMax
	var closestRecord *HitRecord = nil

	for i := 0; i < len(c); i++ {
		if hit, record := c[i].Hit(ray, tMin, closestHit); hit {
			closestRecord = record
			closestHit = record.Distance
			hitAnything = true
		}
	}

	return hitAnything, closestRecord
}

// Bound returns the bounding box of the Collection
func (c Collection) Bound(tMin float64, tMax float64) (bool, *Bbox) {
	if len(c) == 0 {
		return false, nil
	}

	var collectionBox Bbox
	firstBox := true

	for _, actor := range c {
		if bounded, box := actor.shape.Bound(tMin, tMax); bounded {
			if firstBox {
				collectionBox = *box
				firstBox = false
			} else {
				collectionBox = collectionBox.Merge(*box)
			}
		} else {
			return false, nil
		}
	}
	return true, &collectionBox
}

// Comparator returns a comparison function of objects in the collection along the given axis. Used for Index sorting.
func (c Collection) Comparator(startTime, endTime float64, axis int) func(i, j int) bool {
	return func(i, j int) bool {
		leftBound, leftBox := c[i].Bound(startTime, endTime)
		rightBound, rightBox := c[j].Bound(startTime, endTime)

		if !leftBound || !rightBound {
			panic("no bounding box")
		}

		return leftBox.Min.AsArray()[axis] < rightBox.Min.AsArray()[axis]
	}
}

// Index is a binary tree forming a bounding volume hierarchy of objects satisfying the geometry interface
type Index struct {
	box   Bbox
	left  Geometry
	right Geometry
}

// NewIndex builds a bounding volume hierarchy
func NewIndex(world Collection, start, end int, startTime, endTime float64) Index {
	// chose random axis for sorting
	comparator := world.Comparator(startTime, endTime, rand.Intn(3))

	var idx Index
	span := end - start + 1
	if span == 1 {
		idx.left = world[start]
		idx.right = world[start]
	} else if span == 2 {
		if comparator(start, end) {
			idx.left = world[start]
			idx.right = world[start+1]
		} else {
			idx.left = world[start+1]
			idx.right = world[start]
		}
	} else {
		sort.Slice(world, comparator) //TODO does this sort the view or the slice ?
		mid := start + span/2
		idx.left = NewIndex(world, start, mid, startTime, endTime)
		idx.right = NewIndex(world, mid, end, startTime, endTime)
	}

	_, leftBox := idx.left.Bound(startTime, endTime) // TODO error handling if no box
	_, rightBox := idx.right.Bound(startTime, endTime)
	idx.box = leftBox.Merge(*rightBox)
	return idx
}

// Hit implements the hit interface for the Index
func (idx Index) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	if !idx.box.Hit(ray, tMin, tMax) {
		return false, nil
	}

	// two different records to not override by nil pointer if not in right
	hitLeft, recordLeft := idx.left.Hit(ray, tMin, tMax)
	if hitLeft {
		tMax = recordLeft.Distance
	}
	hitRight, recordRight := idx.right.Hit(ray, tMin, tMax)
	if hitRight {
		return true, recordRight
	}
	return hitLeft, recordLeft
}

// Bound returns the bounding box of the Index
func (idx Index) Bound(tMin float64, tMax float64) (bool, *Bbox) {
	return true, &idx.box
}
