package main

import (
	"flag"

	"github.com/mastercactapus/sdf-builder"

	"github.com/deadsy/sdfx/sdf"
)

func init() {
	var (
		shaftDia = 5.2
		// pvcOD   = 42.5
		pvcOD = 21.8
		pvcID = 15.4
	)
	const (
		t = 2
		w = 5
		h = 20

		shaftH = 5
		screwD = 3.5
	)
	flag.Float64Var(&shaftDia, "motor-shaft-dia", shaftDia, "Motor shaft diameter for end cap.")
	flag.Float64Var(&pvcOD, "end-cap-pvc-od", pvcOD, "PVC outer diameter for end caps.")
	flag.Float64Var(&pvcID, "end-cap-pvc-id", pvcID, "PVC inner diameter for end caps.")
	cap := func() builder.Builder {
		return builder.NewCylinder(shaftH, shaftDia+t*2).
			Union(
				builder.NewBox(pvcOD, w, t).RotateZCopy(4),
				builder.NewCylinder(h, pvcOD+t*2).Difference(builder.NewCylinder(h, pvcOD)),
			)

	}
	capB := func() builder.Builder {
		ringDia := pvcID - t*2
		return builder.NewCylinder(shaftH, ringDia).
			Union(
				builder.NewCylinder(3, pvcID+t*2),
				builder.NewCylinder(h, ringDia).Difference(builder.NewCylinder(h, ringDia-t*2)),
				builder.
					NewRoundedBox(t*2, t*2, h+t*2, t).
					SnapMidX(ringDia/2).
					RotateZCopy(3).
					Translate(0, 0, -t).
					Difference(
						builder.NewCylinder(t, pvcID).SnapMaxZ(0),
					),
			)

	}
	Register("motor-cap-b", func() sdf.SDF3 {
		center := capB().
			Difference(
				builder.NewCylinder(shaftH, shaftDia),
			).
			Union(
				builder.NewBox(.5, shaftDia/1.5, shaftH).SnapMaxX(shaftDia / 2),
			)

		return center
	})

	Register("motor-cap", func() sdf.SDF3 {
		center := cap().
			Difference(
				builder.NewCylinder(shaftH, shaftDia),
			).
			Union(
				builder.NewBox(.5, shaftDia/1.5, shaftH).SnapMaxX(shaftDia / 2),
			)

		return center
	})
	Register("idle-cap", func() sdf.SDF3 {
		center := cap().
			Difference(
				builder.NewCylinder(shaftH, screwD),
				builder.NewHexagon(5.5, 2.5).SnapMaxZ(shaftH),
			)

		return center
	})
	Register("idle-cap-b", func() sdf.SDF3 {
		center := capB().
			Difference(
				builder.NewCylinder(shaftH, screwD),
				builder.NewHexagon(5.5, 2.5).SnapMaxZ(shaftH),
			)

		return center
	})

}
