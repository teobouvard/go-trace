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

// computes the location of the hit as "pixel" coordinates for texture mapping
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

// RectXY is a rectangular shape in the XY plane (z=k), bounded by x0, x1, y0 and y1
type RectXY struct {
	x0, x1 float64
	y0, y1 float64
	k      float64
}

// Hit implements the geometry interface for RectXY
func (r RectXY) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	t := (r.k - ray.Origin.Z) / ray.Direction.Z
	if t < tMin || t > tMax {
		return false, nil
	}
	x := ray.Origin.X + t*ray.Direction.X
	y := ray.Origin.Y + t*ray.Direction.Y
	if x < r.x0 || x > r.x1 || y < r.y0 || y > r.y1 {
		return false, nil
	}
	u := (x - r.x0) / (r.x1 - r.x0)
	v := (y - r.y0) / (r.y1 - r.y0)

	// TODO don't forget to check for normal direction in scatter
	return true, &HitRecord{Distance: t, Position: ray.At(t), U: u, V: v, Normal: Vec3{Z: 1}}
}

// Bound returns the bounding box of a RectXY
func (r RectXY) Bound(startTime float64, endTime float64) (bool, *Bbox) {
	return true, &Bbox{Vec3{r.x0, r.y0, r.k - 1e-4}, Vec3{r.x1, r.y1, r.k + 1e-4}}
}

// RectXZ is a rectangular shape in the XZ plane (y=k), bounded by x0, x1, z0 and z1
type RectXZ struct {
	x0, x1 float64
	z0, z1 float64
	k      float64
}

// Hit implements the geometry interface for RectXY
func (r RectXZ) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	t := (r.k - ray.Origin.Y) / ray.Direction.Y
	if t < tMin || t > tMax {
		return false, nil
	}
	x := ray.Origin.X + t*ray.Direction.X
	z := ray.Origin.Z + t*ray.Direction.Z
	if x < r.x0 || x > r.x1 || z < r.z0 || z > r.z1 {
		return false, nil
	}
	u := (x - r.x0) / (r.x1 - r.x0)
	v := (z - r.z0) / (r.z1 - r.z0)

	// TODO don't forget to check for normal direction in scatter
	return true, &HitRecord{Distance: t, Position: ray.At(t), U: u, V: v, Normal: Vec3{Y: 1}}
}

// Bound returns the bounding box of a RectXZ
func (r RectXZ) Bound(startTime float64, endTime float64) (bool, *Bbox) {
	return true, &Bbox{Vec3{r.x0, r.k - 1e-4, r.z0}, Vec3{r.x1, r.k + 1e-4, r.z1}}
}

// RectYZ is a rectangular shape in the YZ plane (x=k), bounded by y0, y1, z0 and z1
type RectYZ struct {
	y0, y1 float64
	z0, z1 float64
	k      float64
}

// Hit implements the geometry interface for RectYZ
func (r RectYZ) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	t := (r.k - ray.Origin.X) / ray.Direction.X
	if t < tMin || t > tMax {
		return false, nil
	}
	y := ray.Origin.Y + t*ray.Direction.Y
	z := ray.Origin.Z + t*ray.Direction.Z
	if y < r.y0 || y > r.y1 || z < r.z0 || z > r.z1 {
		return false, nil
	}
	u := (y - r.y0) / (r.y1 - r.y0)
	v := (z - r.z0) / (r.z1 - r.z0)

	// TODO don't forget to check for normal direction in scatter
	return true, &HitRecord{Distance: t, Position: ray.At(t), U: u, V: v, Normal: Vec3{X: 1}}
}

// Bound returns the bounding box of a RectXZ
func (r RectYZ) Bound(startTime float64, endTime float64) (bool, *Bbox) {
	return true, &Bbox{Vec3{r.k - 1e-4, r.y0, r.z0}, Vec3{r.k + 1e-4, r.y1, r.z1}}
}

// FlipFace is a geometry wrapper for flipping the front face of the wrapped geometry
type FlipFace struct {
	reversed Geometry
}

// Hit returns the hit of the inital geometry, but with the opposed record normal
func (f FlipFace) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	hit, rec := f.reversed.Hit(ray, tMin, tMax)
	if rec != nil {
		rec.Normal = rec.Normal.Scale(-1)
	}
	return hit, rec
}

// Bound returns the bounding box of a the initial geometry
func (f FlipFace) Bound(startTime float64, endTime float64) (bool, *Bbox) {
	return f.reversed.Bound(startTime, endTime)
}

// Box is a cube
type Box struct {
	minPoint Vec3
	maxPoint Vec3
	sides    []Geometry
}

// NewBox constructs a box from its extreme points
func NewBox(p0, p1 Vec3) Box {
	sides := []Geometry{
		RectXY{p0.X, p1.X, p0.Y, p1.Y, p1.Z},
		FlipFace{RectXY{p0.X, p1.X, p0.Y, p1.Y, p0.Z}},
		RectXZ{p0.X, p1.X, p0.Z, p1.Z, p1.Y},
		FlipFace{RectXZ{p0.X, p1.X, p0.Z, p1.Z, p0.Y}},
		RectYZ{p0.Y, p1.Y, p0.Z, p1.Z, p1.X},
		FlipFace{RectYZ{p0.Y, p1.Y, p0.Z, p1.Z, p0.X}},
	}

	return Box{
		minPoint: p0,
		maxPoint: p1,
		sides:    sides,
	}
}

// Hit implements the geometry interface for a Box
func (b Box) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	hitAnything := false
	closestHit := tMax
	var closestRecord *HitRecord = nil

	for _, side := range b.sides {
		if hit, record := side.Hit(ray, tMin, closestHit); hit {
			closestRecord = record
			closestHit = record.Distance
			hitAnything = true
		}
	}

	return hitAnything, closestRecord
}

// Bound returns the bounding box of the Box
func (b Box) Bound(startTime float64, endTime float64) (bool, *Bbox) {
	return true, &Bbox{b.minPoint, b.maxPoint}
}

// Translate is a wrapper around a geometry, which is offset by a translation vector
type Translate struct {
	shape  Geometry
	offset Vec3
}

// Hit implements the geometry interface for a Translated object
// It does so by offsetting the ray rather than the wrapped object
func (t Translate) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	movedRay := Ray{ray.Origin.Sub(t.offset), ray.Direction, ray.Time, ray.RandSource}
	if hit, record := t.shape.Hit(movedRay, tMin, tMax); hit {
		record.Position = record.Position.Add(t.offset)
		return true, record
	}
	return false, nil
}

// Bound returns the bounding box of a translated geometry
func (t Translate) Bound(startTime float64, endTime float64) (bool, *Bbox) {
	if isBounded, bbox := t.shape.Bound(startTime, endTime); isBounded {
		return true, &Bbox{bbox.Min.Add(t.offset), bbox.Max.Add(t.offset)}
	}
	return false, nil
}

// RotateY is a wrapper around a geometry, which is rotated around the Y axis
type RotateY struct {
	shape    Geometry
	sinTheta float64
	cosTheta float64
	bbox     Bbox
	hasBox   bool
}

// NewRotateY constructs a rotated object around the Y axis
func NewRotateY(shape Geometry, angle float64) Geometry {
	theta := angle * math.Pi / 180.0
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)

	hasBox, box := shape.Bound(0, 1) // TODO time should not be guessed here
	minPoint := Vec3{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	maxPoint := Vec3{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				x := float64(i)*box.Max.X + (1.0-float64(i))*box.Min.X
				y := float64(j)*box.Max.Y + (1.0-float64(j))*box.Min.Y
				z := float64(k)*box.Max.Z + (1.0-float64(k))*box.Min.Z

				tmpx := cosTheta*x + sinTheta*z
				tmpz := -sinTheta*x + cosTheta*z

				tmpvec := Vec3{tmpx, y, tmpz}
				minPoint = MinCoord(minPoint, tmpvec)
				maxPoint = MinCoord(maxPoint, tmpvec)
			}
		}
	}

	return RotateY{
		shape:    shape,
		sinTheta: sinTheta,
		cosTheta: cosTheta,
		hasBox:   hasBox,
		bbox:     Bbox{minPoint, maxPoint},
	}
}

// Hit implements the geometry interface for a Rotated object (around Y axis)
func (r RotateY) Hit(ray Ray, tMin float64, tMax float64) (bool, *HitRecord) {
	origin := ray.Origin
	direction := ray.Direction

	origin.X = r.cosTheta*ray.Origin.X - r.sinTheta*ray.Origin.Z
	origin.Z = r.sinTheta*ray.Origin.X + r.cosTheta*ray.Origin.Z

	direction.X = r.cosTheta*ray.Direction.X - r.sinTheta*ray.Direction.Z
	direction.Z = r.sinTheta*ray.Direction.X + r.cosTheta*ray.Direction.Z

	rotatedRay := Ray{origin, direction, ray.Time, ray.RandSource}

	if hit, record := r.shape.Hit(rotatedRay, tMin, tMax); hit {
		pos := record.Position
		n := record.Normal

		pos.X = r.cosTheta*record.Position.X - r.sinTheta*record.Position.Z
		pos.Z = -r.sinTheta*record.Position.X + r.cosTheta*record.Position.Z

		n.X = r.cosTheta*record.Normal.X - r.sinTheta*record.Normal.Z
		n.Z = -r.sinTheta*record.Normal.X + r.cosTheta*record.Normal.Z

		return true, &HitRecord{Distance: record.Distance, Position: pos, Normal: n} // TODO recompute distance ?
	}
	return false, nil
}

// Bound returns the bounding box of a rotated geometry
func (r RotateY) Bound(startTime float64, endTime float64) (bool, *Bbox) {
	return r.hasBox, &r.bbox
}
