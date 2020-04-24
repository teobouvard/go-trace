package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
)

func rayColor(ray Ray, world Collection, depth int) Vec3 {
	if depth <= 0 {
		// too many scattered bounces, assume absorption
		return BLACK
	}

	if hit, record := world.Hit(ray, 0.001, math.MaxFloat64); hit {
		if scatters, attenuation, scattered := record.Material.Scatter(ray, *record); scatters {
			return attenuation.Mul(rayColor(scattered, world, depth-1))
		}
		// texture absorbs all the ray
		return BLACK
	}

	// background white-blue lerp
	unitDirection := ray.Direction.Unit()
	t := 0.5 * (unitDirection.Y + 1.0)
	return WHITE.Scale(1.0 - t).Add(Vec3{0.5, 0.7, 1.0}.Scale(t))
}

func main() {
	rand.Seed(42)
	pixelSamples := 100
	maxScatter := 50

	imageWidth := 200
	imageHeight := 100

	aspectRatio := float64(imageWidth) / float64(imageHeight)
	fov := 60.0
	lookFrom := Vec3{-2.0, 2.0, 1}
	lookAt := Vec3{Z: -1.0}
	up := Vec3{Y: 1}

	camera := NewCamera(lookFrom, lookAt, up, fov, aspectRatio)

	world := Collection{
		[]Actor{
			{
				shape: Sphere{
					Center: Vec3{Y: -100.5, Z: -1},
					Radius: 100,
				},
				texture: Lambertian{
					albedo: Vec3{0.8, 0.8, 0.0},
				},
			},
			{
				shape: Sphere{
					Center: Vec3{Z: -1},
					Radius: 0.5,
				},
				texture: Lambertian{
					albedo: Vec3{0.1, 0.2, 0.5},
				},
			},
			{
				shape: Sphere{
					Center: Vec3{X: 1, Z: -1},
					Radius: 0.5,
				},
				texture: Metal{
					albedo: Vec3{0.8, 0.6, 0.2},
				},
			},
			{
				shape: Sphere{
					Center: Vec3{X: -1, Z: -1},
					Radius: 0.5,
				},
				texture: Dielectric{n: 1.5},
			},
			{
				shape: Sphere{
					Center: Vec3{X: -1, Z: -1},
					Radius: -0.45,
				},
				texture: Dielectric{n: 1.5},
			},
		},
	}

	fmt.Printf("P3\n%v %v\n255\n", imageWidth, imageHeight)
	for j := imageHeight - 1; j >= 0; j-- {
		fmt.Fprintf(os.Stderr, "\rLines remaining: %v", j)
		for i := 0; i < imageWidth; i++ {
			color := WHITE
			for s := 0; s < pixelSamples; s++ {
				u := (float64(i) + rand.Float64()) / float64(imageWidth)
				v := (float64(j) + rand.Float64()) / float64(imageHeight)
				ray := camera.RayTo(u, v)
				color = color.Add(rayColor(ray, world, maxScatter))
			}
			fmt.Printf(color.WriteColor(pixelSamples))
		}
	}
}
