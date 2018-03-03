package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

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
func renderSTL(dir, name string, mesh sdf.SDF3, res float64) {
	if mesh == nil {
		return
	}
	cells := int(mesh.BoundingBox().Size().MaxComponent() / res)
	os.MkdirAll(dir, 0755)
	sdf.RenderSTL(mesh, cells, filepath.Join(dir, name+".stl"))
}

func main() {
	res := flag.Float64("res", 0.2, "Render sampling resolution (smaller values increase detail).")
	assemble := flag.Bool("assemble", false, "Render assembly instructions.")
	dir := flag.String("out", "stl", "Output directory for STL files.")
	flag.Parse()
	if flag.NArg() == 0 {
		log.Println("Processing all", len(names), "objects.")
		for name := range names {
			renderSTL(*dir, name, getMesh(name), *res)
		}
	} else {
		for _, name := range flag.Args() {
			renderSTL(*dir, name, getMesh(name), *res)
		}
	}
	if *assemble {
		renderSTL(*dir, "assemble", assembleMesh(), *res)
	}
}
