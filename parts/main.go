package main

import (
	"flag"
	"log"

	"github.com/deadsy/sdfx/sdf"
)

var names = make(map[string]func() sdf.SDF3)

func Register(name string, fn func() sdf.SDF3) {
	names[name] = fn
}

func getMesh(name string) sdf.SDF3 {
	fn, ok := names[name]
	if !ok {
		log.Println("Not found:", name)
		return nil
	}
	return fn()
}
func renderSTL(name string, res int) {
	mesh := getMesh(name)
	if mesh == nil {
		return
	}
	sdf.RenderSTL(mesh, res, name+".stl")
}
func renderPNG(name string, floor bool) {
	mesh := getMesh(name)
	if mesh == nil {
		return
	}
	sdf.RenderPNG(mesh, floor)
}
func main() {
	resolution := flag.Int("res", 250, "Render sampling resolution (only valid for STL).")
	png := flag.Bool("png", false, "Render PNG instead of STL.")
	floor := flag.Bool("floor", false, "Render the floor (only valid for PNG).")
	flag.Parse()
	if flag.NArg() == 0 {
		log.Println("Processing all", len(names), "objects.")
		for name := range names {
			if *png {
				renderPNG(name, *floor)
			} else {
				renderSTL(name, *resolution)
			}
		}
	} else {
		for _, name := range flag.Args() {
			if *png {
				renderPNG(name, *floor)
			} else {
				renderSTL(name, *resolution)
			}
		}
	}
}
