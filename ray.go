package gotrace

// Ray is a light ray
type Ray struct {
	Origin    Vec3
	Direction Vec3
	Time      float64
}

// At is the point of the ray having travelled t
func (r Ray) At(t float64) Vec3 {
	return r.Origin.Add(r.Direction.Scale(t))
}
