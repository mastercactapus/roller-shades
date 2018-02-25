package main

import (
	"flag"

	"github.com/mastercactapus/sdf-builder"

	"github.com/deadsy/sdfx/sdf"
)

func init() {
	var (
		shaftDia = 5.2
		pvcDia   = 42.5
	)
	const (
		t = 3
		w = 5
		h = 10

		shaftH = 5
	)
	flag.Float64Var(&shaftDia, "motor-cap-shaft-dia", shaftDia, "Shaft diameter for motor cap.")
	flag.Float64Var(&pvcDia, "motor-cap-pvc-dia", pvcDia, "PVC diameter for motor cap.")
	Register("motor-cap", func() sdf.SDF3 {
		center := builder.NewCylinder(shaftH, shaftDia+t*2).
			Union(
				builder.NewBox(pvcDia, w, t).RotateZCopy(4),
			).
			Difference(
				builder.NewCylinder(shaftH, shaftDia),
			).
			Union(
				builder.NewCylinder(h, pvcDia+t*2).Difference(builder.NewCylinder(h, pvcDia)),
				builder.NewBox(.5, shaftDia/2, shaftH).SnapMaxX(shaftDia/2),
			)

		return center
	})
}
