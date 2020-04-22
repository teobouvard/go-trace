package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/teobouvard/gotrace/scene"
	"github.com/teobouvard/gotrace/space"
)

func rayColor(ray scene.Ray, world scene.Collection, depth int) space.Vec3 {
	if depth <= 0 {
		// too many scattered bounces, assume absorption
		return space.NewVec3(0, 0, 0)
	}

	if hit, record := world.Hit(ray, 0.001, math.MaxFloat64); hit {
		if scatters, attenuation, scattered := record.Material().Scatter(ray, record); scatters {
			return space.Mul(attenuation, rayColor(scattered, world, depth-1))
		}
		// material absorbs all the ray
		return space.NewVec3(0, 0, 0)
	}

	// background
	unitDirection := space.Unit(ray.Direction())
	t := 0.5 * (unitDirection.Y() + 1.0)
	white := space.Scale(space.NewVec3(1.0, 1.0, 1.0), 1.0-t)
	blue := space.Scale(space.NewVec3(0.5, 0.7, 1.0), t)
	return space.Add(white, blue)
}

func main() {
	imageWidth := 200
	imageHeight := 100
	pixelSamples := 100
	maxScatter := 50

	fmt.Printf("P3\n%v %v\n255\n", imageWidth, imageHeight)

	camera := scene.NewCamera()

	redish := space.NewVec3(0.7, 0.3, 0.3)
	greenish := space.NewVec3(0.2, 0.8, 0.2)
	smallSphere := scene.NewSphere(space.NewVec3(0, 0, -1), 0.5)
	bigSphere := scene.NewSphere(space.NewVec3(0, -100.5, -1), 100)
	ball := scene.NewActor(smallSphere, scene.NewLambertian(redish))
	earth := scene.NewActor(bigSphere, scene.NewLambertian(greenish))
	world := scene.NewCollection(ball, earth)

	for j := imageHeight - 1; j >= 0; j-- {
		fmt.Fprintf(os.Stderr, "\rLines remaining: %v", j)
		for i := 0; i < imageWidth; i++ {
			color := space.NewVec3(0, 0, 0)
			for s := 0; s < pixelSamples; s++ {
				u := (float64(i) + rand.Float64()) / float64(imageWidth)
				v := (float64(j) + rand.Float64()) / float64(imageHeight)
				ray := camera.LookAt(u, v)
				color = space.Add(color, rayColor(ray, world, maxScatter))

			}
			color.WriteColor(os.Stdout, pixelSamples)
		}
	}
	fmt.Fprintf(os.Stderr, "\n")
}
