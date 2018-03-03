package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	"github.com/mastercactapus/sdf-builder"
)

func init() {
	const (
		thickness = 2
		height    = 7

		diameter = 56

		motorShaftDiameter = 5.2

		slots = 15
	)

	Register("encoder-disc", func() sdf.SDF3 {
		outerRing := builder.NewCylinder(thickness, diameter+thickness*2).
			Difference(builder.NewCylinder(thickness, diameter-thickness*2))

		encodeRing := builder.NewCylinder(height, diameter+thickness).Difference(
			builder.NewCylinder(height, diameter+thickness/2),
		).Difference(
			builder.NewSlice(height, (diameter+thickness)/2, math.Pi/slots).RotateZCopy(slots),
		)

		arms := builder.NewBox(diameter, thickness*2, thickness).RotateZCopy(4)
		shaft := builder.NewCylinder(height-1.5, motorShaftDiameter+thickness*2)

		motorShaft := builder.NewCylinder(height*2, motorShaftDiameter).Difference(
			builder.NewBox(.5, motorShaftDiameter, height*2).Translate(motorShaftDiameter/2-.25, 0, 0),
		)

		mount := arms.Union(shaft).Difference(motorShaft)
		f := sdf.Union3D(
			outerRing,
			encodeRing,
			mount,
		)
		return f
	})
}
