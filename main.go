package main

import (
	"fmt"
	"os"

	"github.com/teobouvard/gotrace/light"
	"github.com/teobouvard/gotrace/space"
)

func hitSphere(center *space.Vec3, radius float64, ray *light.Ray) bool {
	oc := space.Sub(ray.Origin(), center)
	a := space.Dot(ray.Direction(), ray.Direction())
	b := 2.0 * space.Dot(oc, ray.Direction())
	c := space.Dot(oc, oc) - radius*radius
	discriminant := b*b - 4*a*c
	return discriminant > 0
}

func rayColor(ray *light.Ray) *space.Vec3 {
	if hitSphere(space.NewVec3(0, 0, -1), 0.5, ray) {
		return space.NewVec3(1, 0, 0)
	}
	unitDirection := space.Unit(ray.Direction())
	t := 0.5 * (unitDirection.Y() + 1.0)
	c1 := space.Mul(space.NewVec3(1.0, 1.0, 1.0), 1.0-t)
	c2 := space.Mul(space.NewVec3(0.5, 0.7, 1.0), t)
	return space.Add(c1, c2)
}

func main() {
	imageWidth := 200
	imageHeight := 100

	fmt.Printf("P3\n%v %v\n255\n", imageWidth, imageHeight)

	origin := space.NewVec3(0, 0, 0)
	lowerLeft := space.NewVec3(-2, -1, -1)
	horizontal := space.NewVec3(4, 0, 0)
	vertical := space.NewVec3(0, 2, 0)

	for j := imageHeight - 1; j >= 0; j-- {
		fmt.Fprintf(os.Stderr, "\rLines remaining: %v", j)
		for i := 0; i < imageWidth; i++ {
			u := float64(i) / float64(imageWidth)
			v := float64(j) / float64(imageHeight)
			dir := space.Add(lowerLeft, space.Mul(horizontal, u), space.Mul(vertical, v))
			r := light.NewRay(origin, dir)
			color := rayColor(r)
			color.WriteColor(os.Stdout)
		}
	}
	fmt.Fprintf(os.Stderr, "\n")
}
