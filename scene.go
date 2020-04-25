package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"

	"golang.org/x/sync/semaphore"
)

// Scene is the whole scene to be rendered
type Scene struct {
	world        Collection
	camera       Camera
	pixelSamples int
	width        int
	height       int
	maxScatter   int
}

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

// Render renders the scene
func (s Scene) Render() {
	// allocate memory for result
	lines := make([][]string, s.height)
	for i := range lines {
		lines[i] = make([]string, s.width)
	}

	// create workgroup
	ctx := context.TODO()
	nWorkers := int64(runtime.NumCPU())
	sem := semaphore.NewWeighted(nWorkers)
	for j := 0; j < s.height; j++ {
		sem.Acquire(ctx, 1)
		go func(j int) {
			defer sem.Release(1)
			fmt.Fprintf(os.Stderr, "\rLines remaining: %v", s.height-j)
			for i := 0; i < s.width; i++ {
				color := WHITE
				for k := 0; k < s.pixelSamples; k++ {
					u := (float64(i) + rand.Float64()) / float64(s.width)
					v := (float64(j) + rand.Float64()) / float64(s.height)
					ray := s.camera.RayTo(u, v)
					color = color.Add(rayColor(ray, s.world, s.maxScatter))
				}
				lines[s.height-j-1][i] = fmt.Sprintf(color.WriteColor(s.pixelSamples))
			}
		}(j)
	}

	// wait for all workers to exit
	sem.Acquire(ctx, nWorkers)

	// write image
	fmt.Printf("P3\n%v %v\n255\n", s.width, s.height)
	for _, line := range lines {
		for _, col := range line {
			fmt.Print(col)
		}
	}
}

/* BOOK SCENE */

// BookScene creates the scene on the cover of the first book
func BookScene() Scene {
	// image settings
	imageWidth := 1440
	imageHeight := 1080
	pixelSamples := 100
	maxScatter := 100

	// camera settings
	aspectRatio := float64(imageWidth) / float64(imageHeight)
	fov := 20.0
	lookFrom := Vec3{13, 2, 3}
	lookAt := Vec3{0, 0, 0}
	up := Vec3{Y: 1}
	focusDist := 10.0
	aperture := 0.1
	camera := NewCamera(lookFrom, lookAt, up, fov, aspectRatio, aperture, focusDist)

	// objects on the scene
	world := Collection{
		[]Actor{
			{
				shape: Sphere{
					Center: Vec3{X: 0, Y: -1000, Z: 0},
					Radius: 1000,
				},
				texture: Lambertian{
					albedo: Vec3{0.5, 0.5, 0.5},
				},
			},
		},
	}

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := Vec3{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}
			randMaterial := rand.Float64()
			noBalls := Vec3{4, 0.2, 0}
			if center.Sub(noBalls).Norm() > 0.9 {
				if randMaterial < 0.8 {
					// diffuse
					albedo := RandVec().Mul(RandVec())
					actor := Actor{
						shape: Sphere{
							Center: center,
							Radius: 0.2,
						},
						texture: Lambertian{
							albedo: albedo,
						},
					}
					world.Add(actor)
				} else if randMaterial < 0.95 {
					//metal
					albedo := RandVecInterval(0.5, 1.0)
					fuzz := rand.Float64() / 2
					actor := Actor{
						shape: Sphere{
							Center: center,
							Radius: 0.2,
						},
						texture: Metal{
							albedo: albedo,
							fuzz:   fuzz,
						},
					}
					world.Add(actor)
				} else {
					// glass
					actor := Actor{
						shape: Sphere{
							Center: center,
							Radius: 0.2,
						},
						texture: Dielectric{
							n: 1.5,
						},
					}
					world.Add(actor)
				}
			}
		}
	}

	world.Add(
		Actor{
			shape: Sphere{
				Center: Vec3{Y: 1},
				Radius: 1,
			},
			texture: Dielectric{
				n: 1.5,
			},
		},
		Actor{
			shape: Sphere{
				Center: Vec3{X: -4, Y: 1},
				Radius: 1,
			},
			texture: Lambertian{
				albedo: Vec3{0.4, 0.2, 0.1},
			},
		},
		Actor{
			shape: Sphere{
				Center: Vec3{X: 4, Y: 1},
				Radius: 1,
			},
			texture: Metal{
				albedo: Vec3{0.7, 0.6, 0.5},
				fuzz:   0,
			},
		},
	)

	return Scene{world, camera, pixelSamples, imageWidth, imageHeight, maxScatter}
}
