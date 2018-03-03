package main

import (
	"flag"
	"log"

	"github.com/deadsy/sdfx/sdf"
)

var names = make(map[string]func() sdf.SDF3)
var cache = make(map[string]sdf.SDF3)

// Register will add a new part to be available for rendering.
func Register(name string, fn func() sdf.SDF3) {
	names[name] = fn
}

func getMesh(name string) sdf.SDF3 {
	if s, ok := cache[name]; ok {
		return s
	}
	fn, ok := names[name]
	if !ok {
		log.Println("Not found:", name)
		return nil
	}
	cache[name] = fn()
	return cache[name]
}
func renderSTL(name string, mesh sdf.SDF3, ppmm float64) {
	if mesh == nil {
		return
	}
	res := int(mesh.BoundingBox().Size().MaxComponent() * ppmm)
	sdf.RenderSTL(mesh, res, name+".stl")
}

func main() {
	resolution := flag.Float64("res", 5, "Render sampling resolution (pixels per mm).")
	assemble := flag.Bool("assemble", false, "Render assembly instructions.")
	flag.Parse()
	if flag.NArg() == 0 {
		log.Println("Processing all", len(names), "objects.")
		for name := range names {
			renderSTL(name, getMesh(name), *resolution)
		}
	} else {
		for _, name := range flag.Args() {
			renderSTL(name, getMesh(name), *resolution)
		}
	}
	if *assemble {
		renderSTL("assemble", assembleMesh(), *resolution)
	}
}
