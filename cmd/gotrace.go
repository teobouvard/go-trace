package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"runtime/pprof"

	"github.com/teobouvard/gotrace"
)

var (
	// TODO add quality options and scene parsing from file
	cpuProfile  = flag.String("profile", "perf", "write cpu profile to file")
	outputImage = flag.String("output", "render.png", "output rendered image to file")
)

func main() {
	flag.Parse()

	// execution profiling
	// use go tool pprof perf, and web/topX
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *outputImage != "" {
		if _, err := os.Stat(*outputImage); os.IsExist(err) {
			fmt.Println("Output file already exists")
		} else {
			f, err := os.Create(*outputImage)
			if err != nil {
				log.Fatal(err)
			}
			scene := gotrace.FinalScene()
			img := scene.Render(2000, -1, 5000, 100)
			//img := scene.Render(300, -1, 50, 50)
			png.Encode(f, img)
		}
	}

}
