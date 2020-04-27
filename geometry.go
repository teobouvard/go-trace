package gotrace

import (
	"math"

	"github.com/teobouvard/gotrace/util"
)

/*
Geometry interface

Hit

@in
	ray : a light ray
	tMin : closer objects are not considered
	tMax : further objects are not considered
@out
	bool : if the ray hit the geometry
	HitRecord : information about the hit, or nil

Bound

@in
	startTime : the starting time for bounding
	endTime : the ending time for bounding
@out
	bool : if the geometry can be bounded (false for infinite planes)
	Bbox : bounding box (aabb) of the geometry, if applicable
*/
type Geometry interface {
	Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord)
	Bound(startTime float64, endTime float64) (bool, *Bbox)
}

// Sphere geometry
type Sphere struct {
	Center Vec3
	Radius float64
}

// computes the location of the hit as "pixel" coordinates
func (s Sphere) pixelHit(pos Vec3) (u, v float64) {
	phi := math.Atan2(pos.Z, pos.X)
	theta := math.Asin(pos.Y)
	u = 1 - (phi+math.Pi)/(2*math.Pi)
	v = (theta + math.Pi/2) / math.Pi
	return
}

// Hit implements the geomtry interface for checking the intersection of a Ray and a Sphere
func (s Sphere) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	oc := ray.Origin.Sub(s.Center)
	a := ray.Direction.SquareNorm()
	b := oc.Dot(ray.Direction)
	c := oc.SquareNorm() - s.Radius*s.Radius
	discriminant := b*b - a*c

	if discriminant > 0 {
		root := math.Sqrt(discriminant)
		// first quadratic solution, closest to camera
		t := (-b - root) / a
		if t < tMax && t > tMin {
			pos := ray.At(t)
			/*
				Previously, I thought doing pos.Sub(s.Center).Unit() was smarter than to divide by the radius.
				This led to a very nasty bug when using negative radii as the normal was computed on the wrong side of the geometry.
			*/
			n := pos.Sub(s.Center).Div(s.Radius)
			u, v := s.pixelHit(n)
			return true, &HitRecord{Distance: t, Position: pos, Normal: n, U: u, V: v}
		}
		// second solution, farthest from camera
		t = (-b + root) / a
		if t < tMax && t > tMin {
			pos := ray.At(t)
			n := pos.Sub(s.Center).Div(s.Radius)
			u, v := s.pixelHit(n)
			return true, &HitRecord{Distance: t, Position: pos, Normal: n, U: u, V: v}
		}
	}

	return false, nil
}

// Bound returns the bounding box of the Sphere
func (s Sphere) Bound(startTime float64, endTime float64) (bool, *Bbox) {
	bounds := Vec3{s.Radius, s.Radius, s.Radius}
	box := Bbox{
		Min: s.Center.Sub(bounds),
		Max: s.Center.Add(bounds),
	}
	return true, &box
}

// MovingSphere geometry
type MovingSphere struct {
	CenterStart   Vec3
	CenterStop    Vec3
	Radius        float64
	tStart, tStop float64
}

func (s MovingSphere) centerAt(time float64) Vec3 {
	elapsed := util.Map(time, s.tStart, s.tStop, 0, 1)
	moved := s.CenterStop.Sub(s.CenterStart).Scale(elapsed)
	return s.CenterStart.Add(moved)
}

// Hit implements the geomtry interface for checking the intersection of a Ray and a MovingSphere
func (s MovingSphere) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	center := s.centerAt(ray.Time)
	oc := ray.Origin.Sub(center)
	a := ray.Direction.SquareNorm()
	b := oc.Dot(ray.Direction)
	c := oc.SquareNorm() - s.Radius*s.Radius
	discriminant := b*b - a*c

	if discriminant > 0 {
		root := math.Sqrt(discriminant)
		// first solution, closest to camera
		t := (-b - root) / a
		if t < tMax && t > tMin {
			pos := ray.At(t)
			n := pos.Sub(center).Div(s.Radius)
			return true, &HitRecord{Distance: t, Position: pos, Normal: n}
		}
		// second solution, farthest from camera
		t = (-b + root) / a
		if t < tMax && t > tMin {
			pos := ray.At(t)
			n := pos.Sub(center).Div(s.Radius)
			return true, &HitRecord{Distance: t, Position: pos, Normal: n}
		}
	}

	return false, nil
}

// Bound returns the bounding box of the MovingSphere
func (s MovingSphere) Bound(startTime float64, endTime float64) (bool, *Bbox) {
	bounds := Vec3{s.Radius, s.Radius, s.Radius}
	startBox := Bbox{
		Min: s.centerAt(startTime).Sub(bounds),
		Max: s.centerAt(startTime).Add(bounds),
	}
	stopBox := Bbox{
		Min: s.centerAt(endTime).Sub(bounds),
		Max: s.centerAt(endTime).Add(bounds),
	}
	box := startBox.Merge(stopBox)
	return true, &box
}
