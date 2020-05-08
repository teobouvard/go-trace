package gotrace

import (
	"context"
	"image"
	"image/draw"
	_ "image/jpeg"
	"log"
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
	world      *Index
	camera     Camera
	background Vec3
}

// NewScene creates a scene that can be rendered. It contains all actors in the world collection, and is viewed from the camera.
func NewScene(camera Camera, world Collection, background Vec3) *Scene {
	return &Scene{
		world:      NewIndex(world, 0, len(world)-1, camera.tStart, camera.tStop),
		camera:     camera,
		background: background,
	}

}

func (s *Scene) rayColor(ray Ray, depth int) Vec3 {
	if depth <= 0 {
		// too many scattered bounces, assume absorption
		return BLACK
	}

	if hit, record := s.world.Hit(ray, 0.001, math.MaxFloat64); hit {
		emitted := record.Material.Emit(record.U, record.V, record.Position)
		if scatters, attenuation, scattered := record.Material.Scatter(ray, *record); scatters {
			return emitted.Add(attenuation.Mul(s.rayColor(scattered, depth-1)))
		}
		return emitted
	}

	return s.background
}

// Render renders the scene with the given parameters
func (s *Scene) Render(width, height, pixelSamples, maxScatter int) *image.RGBA {
	// set default value for max number of ray bounces before absorption
	if maxScatter <= 0 {
		maxScatter = 50
	}

	if pixelSamples <= 0 {
		pixelSamples = 50
	}

	// deduce height from aspect ratio
	if height == -1 {
		height = int(float64(width) / s.camera.AspectRatio)
	}

	// create image
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// create workgroup to render one line per available thread
	ctx := context.TODO()
	nWorkers := int64(runtime.NumCPU())
	sem := semaphore.NewWeighted(nWorkers)
	bar := pb.StartNew(height)
	for j := 0; j < height; j++ {
		sem.Acquire(ctx, 1) // why not in goroutine ?
		go func(j int) {
			defer sem.Release(1)
			rnd := rand.New(rand.NewSource(int64(42 * j)))
			for i := 0; i < width; i++ {
				pixel := BLACK
				for k := 0; k < pixelSamples; k++ {
					u := (float64(i) + rnd.Float64()) / float64(width)
					v := (float64(j) + rnd.Float64()) / float64(height)
					ray := s.camera.RayTo(u, v, rnd)
					pixel = pixel.Add(s.rayColor(ray, maxScatter))
				}
				// set image color
				imgColor := pixel.GetColor(pixelSamples)
				// why isn't there a data race if access to the image is not protected ?
				img.Set(i, height-j-1, imgColor)
			}
			bar.Increment()
		}(j)
	}

	// wait for all workers to exit
	sem.Acquire(ctx, nWorkers)
	bar.Finish()
	return img
}

// BookScene creates the scene on the cover of the first book
func BookScene() *Scene {
	// camera settings
	aspectRatio := 2.0
	fov := 20.0
	lookFrom := Vec3{13, 2, 3}
	lookAt := Vec3{0, 0, 0}
	up := Vec3{Y: 1}
	focusDist := 10.0
	aperture := 0.1
	camera := NewCamera(lookFrom, lookAt, up, fov, aspectRatio, aperture, focusDist, 0, 1)

	rnd := rand.New(rand.NewSource(42))
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
					albedo := RandVec(rnd).Mul(RandVec(rnd))
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
					albedo := RandVecInterval(0.5, 1.0, rnd)
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

	return NewScene(camera, objects, WHITE)
}

// MovingSpheres creates the scene on the cover of the first book, with bouncing balls
func MovingSpheres() *Scene {
	// camera settings
	aspectRatio := 2.0
	fov := 20.0
	lookFrom := Vec3{13, 2, 3}
	lookAt := Vec3{0, 0, 0}
	up := Vec3{Y: 1}
	focusDist := 10.0
	aperture := 0.0
	startTime := 0.0
	endTime := 1.0
	camera := NewCamera(lookFrom, lookAt, up, fov, aspectRatio, aperture, focusDist, startTime, endTime)

	// objects on the scene
	rnd := rand.New(rand.NewSource(42))
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
					freq: 10,
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
					albedo := RandVec(rnd).Mul(RandVec(rnd))
					actor := Actor{
						shape: MovingSphere{
							CenterStart: center,
							CenterStop:  center.Add(Vec3{Y: rand.Float64() / 2.0}),
							tStart:      startTime,
							tStop:       endTime,
							Radius:      0.2,
						},
						material: Lambertian{
							albedo: ConstantTexture{albedo},
						},
					}
					world.Add(actor)
				} else if randMaterial < 0.95 {
					//metal
					albedo := RandVecInterval(0.5, 1.0, rnd)
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

	return NewScene(camera, world, WHITE)
}

// MarbleScene is a scene with a black and white marble
func MarbleScene() *Scene {
	aspectRatio := 16.0 / 9.0
	fov := 33.0
	lookFrom := Vec3{13, 2, 3}
	lookAt := Vec3{0, 2, 0}
	up := Vec3{Y: 1}
	focusDist := 10.0
	aperture := 0.0
	camera := NewCamera(lookFrom, lookAt, up, fov, aspectRatio, aperture, focusDist, 0, 1)

	objects := Collection{
		Actor{
			shape: Sphere{
				Center: Vec3{X: 5, Y: 2, Z: 3},
				Radius: 2,
			},
			material: Lambertian{
				Marble{
					noise:      opensimplex.New(51),
					depth:      7,
					turbulence: 5,
					scale:      4,
				},
			},
		},
	}

	return NewScene(camera, objects, WHITE)
}

// EarthScene is a scene demonstrating image textures
func EarthScene() *Scene {
	// camera settings
	aspectRatio := 2.0
	fov := 33.0
	lookFrom := Vec3{13, 2, 3}
	lookAt := Vec3{0, 2, 0}
	up := Vec3{Y: 1}
	focusDist := 10.0
	aperture := 0.0
	camera := NewCamera(lookFrom, lookAt, up, fov, aspectRatio, aperture, focusDist, 0, 1)

	f, err := os.Open("./assets/blue_marble_september.jpg")
	src, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	bounds := src.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(img, img.Bounds(), src, bounds.Min, draw.Src)

	objects := Collection{
		Actor{
			shape: Sphere{
				Center: Vec3{Y: -1000},
				Radius: 1000,
			},
			material: Lambertian{
				albedo: ConstantTexture{Vec3{0.4, 0.4, 0.4}},
			},
		},
		Actor{
			shape: Sphere{
				Center: Vec3{X: 5, Y: 2, Z: 3},
				Radius: 2,
			},
			material: Lambertian{
				Image{
					data: img,
				},
			},
		},
	}

	return NewScene(camera, objects, WHITE)
}

// LightMarbleScene is a scene with a black and white marble with lights
func LightMarbleScene() *Scene {
	// camera settings
	aspectRatio := 2.0
	fov := 50.0
	lookFrom := Vec3{13, 3, 3}
	lookAt := Vec3{0, 2, 0}
	up := Vec3{Y: 1}
	focusDist := 10.0
	aperture := 0.0
	camera := NewCamera(lookFrom, lookAt, up, fov, aspectRatio, aperture, focusDist, 0, 1)

	objects := Collection{
		Actor{
			shape: Sphere{
				Center: Vec3{Y: -1000},
				Radius: 1000,
			},
			material: Lambertian{
				albedo: ConstantTexture{Vec3{0.8, 0.8, 0.8}},
			},
		},
		Actor{
			shape: Sphere{
				Center: Vec3{X: 5, Y: 2, Z: 3},
				Radius: 2,
			},
			material: Lambertian{
				Marble{
					noise:      opensimplex.New(51),
					depth:      7,
					turbulence: 5,
					scale:      4,
				},
			},
		},
		Actor{
			shape: Sphere{
				Center: Vec3{X: 7, Y: 4, Z: -1},
				Radius: 1,
			},
			material: DiffuseLight{
				ConstantTexture{WHITE.Scale(5)},
			},
		},
	}

	return NewScene(camera, objects, WHITE)
}

// CornellBox is a the classic cornell box scene
func CornellBox() *Scene {
	// camera settings
	aspectRatio := 1.0
	fov := 40.0
	lookFrom := Vec3{278, 278, -800}
	lookAt := Vec3{278, 278, 0}
	up := Vec3{Y: 1}
	focusDist := 10.0
	aperture := 0.0
	camera := NewCamera(lookFrom, lookAt, up, fov, aspectRatio, aperture, focusDist, 0, 1)

	objects := Collection{
		// left wall - green
		Actor{
			shape: FlipFace{RectYZ{0, 555, 0, 555, 555}},
			material: Lambertian{
				albedo: ConstantTexture{Vec3{0.12, 0.45, 0.15}},
			},
		},
		// right wall - red
		Actor{
			shape: RectYZ{0, 555, 0, 555, 0},
			material: Lambertian{
				albedo: ConstantTexture{Vec3{0.65, 0.05, 0.05}},
			},
		},
		// roof light
		Actor{
			shape: RectXZ{213, 343, 227, 332, 554},
			material: DiffuseLight{
				emit: ConstantTexture{WHITE.Scale(15)},
			},
		},
		// floor
		Actor{
			shape: RectXZ{0, 555, 0, 555, 0},
			material: Lambertian{
				albedo: ConstantTexture{WHITE},
			},
		},
		// ceiling
		Actor{
			shape: FlipFace{RectXZ{0, 555, 0, 555, 555}},
			material: Lambertian{
				albedo: ConstantTexture{WHITE},
			},
		},
		// back wall
		Actor{
			shape: FlipFace{RectXY{0, 555, 0, 555, 555}},
			material: Lambertian{
				albedo: ConstantTexture{WHITE},
			},
		},
		// Front box
		Actor{
			shape: NewBox(Vec3{130, 0, 65}, Vec3{295, 165, 230}),
			material: Lambertian{
				albedo: ConstantTexture{WHITE},
			},
		},
		// Back box
		Actor{
			shape: NewBox(Vec3{265, 0, 295}, Vec3{430, 330, 460}),
			material: Lambertian{
				albedo: ConstantTexture{WHITE},
			},
		},
	}
	return NewScene(camera, objects, BLACK)
}
