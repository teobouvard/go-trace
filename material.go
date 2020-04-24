package main

// Material define the way actors interact with a ray
type Material interface {
	Scatter(ray Ray, hit HitRecord) (bool, Vec3, Ray)
}
