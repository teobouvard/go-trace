package main

// Ray is a light ray
type Ray struct {
	Origin    Vec3
	Direction Vec3
}

// At is the point of the ray after it travelled t units of time
func (r Ray) At(t float64) Vec3 {
	return r.Origin.Add(r.Direction.Scale(t))
}
