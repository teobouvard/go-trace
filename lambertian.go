package main

// Lambertian is a diffuse texture
type Lambertian struct {
	albedo Vec3
}

// Scatter defines how a lambertian material scatters a Ray
func (l Lambertian) Scatter(ray Ray, hit HitRecord) (bool, Vec3, Ray) {
	scatterDirection := hit.Normal.Add(RandLambertian())
	scattered := Ray{hit.Position, scatterDirection}
	return true, l.albedo, scattered
}
