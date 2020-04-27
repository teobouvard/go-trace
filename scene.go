package gotrace

import (
	"context"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"runtime"

	"github.com/cheggaaa/pb/v3"
	"github.com/ojrac/opensimplex-go"
	"golang.org/x/sync/semaphore"
)

// Scene is the whole scene to be rendered
type Scene struct {
	world        Index
	camera       Camera
	pixelSamples int
	width        int
	height       int
	maxScatter   int
}

func (s Scene) rayColor(ray Ray, depth int) Vec3 {
	if depth <= 0 {
		// too many scattered bounces, assume absorption
		return BLACK
	}

	if hit, record := s.world.Hit(ray, 0.001, math.MaxFloat64); hit {
		if scatters, attenuation, scattered := record.Material.Scatter(ray, *record); scatters {
			return attenuation.Mul(s.rayColor(scattered, depth-1))
		}
		// material absorbs all the ray
		return BLACK
	}

	// background white-blue lerp
	unitDirection := ray.Direction.Unit()
	t := 0.5 * (unitDirection.Y + 1.0)
	return WHITE.Scale(1.0 - t).Add(Vec3{0.5, 0.7, 1.0}.Scale(t))
}

// Render renders the scene
func (s Scene) Render() {
	// create image
	upLeft := image.Point{0, 0}
	lowRight := image.Point{s.width, s.height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// create workgroup
	ctx := context.TODO()
	nWorkers := int64(runtime.NumCPU())
	sem := semaphore.NewWeighted(nWorkers)
	bar := pb.StartNew(s.height)
	for j := 0; j < s.height; j++ {
		sem.Acquire(ctx, 1)
		go func(j int) {
			defer sem.Release(1)
			for i := 0; i < s.width; i++ {
				pixel := BLACK
				for k := 0; k < s.pixelSamples; k++ {
					u := (float64(i) + rand.Float64()) / float64(s.width)
					v := (float64(j) + rand.Float64()) / float64(s.height)
					ray := s.camera.RayTo(u, v)
					pixel = pixel.Add(s.rayColor(ray, s.maxScatter))
				}
				// set image color
				imgColor := pixel.GetColor(s.pixelSamples)
				img.Set(i, s.height-j-1, imgColor)
			}
			bar.Increment()
		}(j)
	}

	// wait for all workers to exit
	sem.Acquire(ctx, nWorkers)

	// write image
	os.Remove("img/image.png")
	f, _ := os.Create("img/image.png")
	png.Encode(f, img)
}

// BookScene creates the scene on the cover of the first book
func BookScene() Scene {
	// image settings
	//imageWidth := 1440
	//imageHeight := 1080
	imageWidth := 200
	imageHeight := 100
	pixelSamples := 100
	maxScatter := 50

	// camera settings
	aspectRatio := float64(imageWidth) / float64(imageHeight)
	fov := 20.0
	lookFrom := Vec3{13, 2, 3}
	lookAt := Vec3{0, 0, 0}
	up := Vec3{Y: 1}
	focusDist := 10.0
	aperture := 0.1
	camera := NewCamera(lookFrom, lookAt, up, fov, aspectRatio, aperture, focusDist, 0, 1)

	// objects on the scene
	objects := Collection{
		{
			shape: Sphere{
				Center: Vec3{X: 0, Y: -1000, Z: 0},
				Radius: 1000,
			},
			material: Lambertian{
				albedo: ConstantTexture{Vec3{0.5, 0.5, 0.5}},
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
						material: Lambertian{
							albedo: ConstantTexture{albedo},
						},
					}
					objects.Add(actor)
				} else if randMaterial < 0.95 {
					//metal
					albedo := RandVecInterval(0.5, 1.0)
					fuzz := rand.Float64() / 2
					actor := Actor{
						shape: Sphere{
							Center: center,
							Radius: 0.2,
						},
						material: Metal{
							albedo: albedo,
							fuzz:   fuzz,
						},
					}
					objects.Add(actor)
				} else {
					// glass
					actor := Actor{
						shape: Sphere{
							Center: center,
							Radius: 0.2,
						},
						material: Dielectric{
							n: 1.5,
						},
					}
					objects.Add(actor)
				}
			}
		}
	}

	objects.Add(
		Actor{
			shape: Sphere{
				Center: Vec3{Y: 1},
				Radius: 1,
			},
			material: Dielectric{
				n: 1.5,
			},
		},
		Actor{
			shape: Sphere{
				Center: Vec3{X: -4, Y: 1},
				Radius: 1,
			},
			material: Lambertian{
				albedo: ConstantTexture{Vec3{0.4, 0.2, 0.1}},
			},
		},
		Actor{
			shape: Sphere{
				Center: Vec3{X: 4, Y: 1},
				Radius: 1,
			},
			material: Metal{
				albedo: Vec3{0.7, 0.6, 0.5},
				fuzz:   0,
			},
		},
	)

	world := NewIndex(objects, 0, len(objects)-1, 0, 1)
	return Scene{world, camera, pixelSamples, imageWidth, imageHeight, maxScatter}
}

// MovingSpheres creates the scene on the cover of the first book, with bouncing balls
func MovingSpheres() Scene {
	// image settings
	//imageWidth := 1440
	//imageHeight := 1080
	imageWidth := 200
	imageHeight := 100
	pixelSamples := 100
	maxScatter := 50

	// camera settings
	aspectRatio := float64(imageWidth) / float64(imageHeight)
	fov := 20.0
	lookFrom := Vec3{13, 2, 3}
	lookAt := Vec3{0, 0, 0}
	up := Vec3{Y: 1}
	focusDist := 10.0
	aperture := 0.0
	camera := NewCamera(lookFrom, lookAt, up, fov, aspectRatio, aperture, focusDist, 0, 1)

	// objects on the scene
	world := Collection{
		{
			shape: Sphere{
				Center: Vec3{X: 0, Y: -1000, Z: 0},
				Radius: 1000,
			},
			material: Lambertian{
				albedo: CheckerTexture{
					odd:  ConstantTexture{Vec3{0, 0, 0}},
					even: ConstantTexture{Vec3{1, 1, 1}},
				},
			},
		},
	}

	for a := -10; a < 10; a++ {
		for b := -10; b < 10; b++ {
			center := Vec3{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + 0.9*rand.Float64()}
			randMaterial := rand.Float64()
			noBalls := Vec3{4, 0.2, 0}
			if center.Sub(noBalls).Norm() > 0.9 {
				if randMaterial < 0.8 {
					// diffuse
					albedo := RandVec().Mul(RandVec())
					actor := Actor{
						shape: MovingSphere{
							CenterStart: center,
							CenterStop:  center.Add(Vec3{Y: rand.Float64() / 2.0}),
							tStart:      0,
							tStop:       1,
							Radius:      0.2,
						},
						material: Lambertian{
							albedo: ConstantTexture{albedo},
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
						material: Metal{
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
						material: Dielectric{
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
			material: Dielectric{
				n: 1.5,
			},
		},
		Actor{
			shape: Sphere{
				Center: Vec3{X: -4, Y: 1},
				Radius: 1,
			},
			material: Lambertian{
				albedo: ConstantTexture{Vec3{0.4, 0.2, 0.1}},
			},
		},
		Actor{
			shape: Sphere{
				Center: Vec3{X: 4, Y: 1},
				Radius: 1,
			},
			material: Metal{
				albedo: Vec3{0.7, 0.6, 0.5},
				fuzz:   0,
			},
		},
	)
	index := NewIndex(world, 0, len(world)-1, 0, 1)
	return Scene{index, camera, pixelSamples, imageWidth, imageHeight, maxScatter}
}

// NoisyScene is a scene with Opensimplex noise
func NoisyScene() Scene {
	imageWidth := 200
	imageHeight := 100
	pixelSamples := 100
	maxScatter := 50

	// camera settings
	aspectRatio := float64(imageWidth) / float64(imageHeight)
	fov := 20.0
	lookFrom := Vec3{13, 2, 3}
	lookAt := Vec3{0, 0, 0}
	up := Vec3{Y: 1}
	focusDist := 10.0
	aperture := 0.0
	camera := NewCamera(lookFrom, lookAt, up, fov, aspectRatio, aperture, focusDist, 0, 1)

	noise := Noise{opensimplex.New(42)}
	objects := Collection{
		Actor{
			shape: Sphere{
				Center: Vec3{Y:-1000},
				Radius: 1000,
			},
			material: Lambertian{noise},
		},
	}
	world := NewIndex(objects, 0, len(objects)-1, 0, 1)
	return Scene{world, camera, pixelSamples, imageWidth, imageHeight, maxScatter}
}