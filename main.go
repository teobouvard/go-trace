package main

import (
	"fmt"
	"os"

	"github.com/teobouvard/go-trace/space"
)

func main() {
	imageWidth := 400
	imageHeight := 200

	fmt.Printf("P3\n%v %v\n255\n", imageWidth, imageHeight)

	for j := imageHeight - 1; j >= 0; j-- {
		fmt.Fprintf(os.Stderr, "\rLines remaining: %v", j)
		for i := 0; i < imageWidth; i++ {
			v := space.NewVec3(float64(i)/float64(imageWidth), float64(j)/float64(imageHeight), 0.2)
			v.WriteColor(os.Stdout)
		}
	}
	fmt.Fprintf(os.Stderr, "\n")
}
