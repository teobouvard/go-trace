package main

import "math"

/*
Geometry interface
@in
	ray : a light ray
	tMin : closer objects are not considered
	tMax : further objects are not considered
@out
	bool : if the ray hit the geometry
	float64 : t at hit
	Vec3 : position of hit
	Vec3 : normal of geometry at hit
*/
type Geometry interface {
	Hit(ray Ray, tMin float64, tMax float64) (bool, float64, Vec3, Vec3)
}

// Sphere geometry
type Sphere struct {
	Center Vec3
	Radius float64
}

// Hit implements the geomtry interface for checking the intersection of a Ray and a Sphere
func (s Sphere) Hit(ray Ray, tMin float64, tMax float64) (bool, float64, Vec3, Vec3) {
	oc := ray.Origin.Sub(s.Center)
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
			/*
				Previously, I thought doing pos.Sub(s.Center).Unit() was smarter than to divide by the radius.
				This led to a very nasty bug when using negative radii as the normal was computed on the wrong side of the geometry.
			*/
			normal := pos.Sub(s.Center).Div(s.Radius)
			return true, t, pos, normal
		}
		// second solution, farthest from camera
		t = (-b + root) / a
		if t < tMax && t > tMin {
			pos := ray.At(t)
			normal := pos.Sub(s.Center).Div(s.Radius)
			return true, t, pos, normal
		}
	}

	return false, -1, Vec3{}, Vec3{}
}

// MovingSphere geometry
type MovingSphere struct {
	CenterStart   Vec3
	CenterStop    Vec3
	Radius        float64
	tStart, tStop float64
}

func (s MovingSphere) centerAt(time float64) Vec3 {
	elapsed := (time - s.tStart) / (s.tStop - s.tStart)
	moved := s.CenterStop.Sub(s.CenterStart).Scale(elapsed)
	return s.CenterStart.Add(moved)
}

// Hit implements the geomtry interface for checking the intersection of a Ray and a MovingSphere
func (s MovingSphere) Hit(ray Ray, tMin float64, tMax float64) (bool, float64, Vec3, Vec3) {
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
			/*
				Previously, I thought doing pos.Sub(s.Center).Unit() was smarter than to divide by the radius.
				This led to a very nasty bug when using negative radii as the normal was computed on the wrong side of the geometry.
			*/
			normal := pos.Sub(center).Div(s.Radius)
			return true, t, pos, normal
		}
		// second solution, farthest from camera
		t = (-b + root) / a
		if t < tMax && t > tMin {
			pos := ray.At(t)
			normal := pos.Sub(center).Div(s.Radius)
			return true, t, pos, normal
		}
	}

	return false, -1, Vec3{}, Vec3{}
}
