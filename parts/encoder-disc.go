package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	. "github.com/mastercactapus/sdf-builder"
)

func Slice(height, length, angle float64) Builder {
	o := math.Tan(angle/2) * length
	return Builder{
		SDF3: sdf.Extrude3D(sdf.Polygon2D([]sdf.V2{
			{},
			{X: length, Y: o},
			{X: length, Y: -o},
		}), height),
	}.Translate(0, 0, height/2)
}
func init() {
	const (
		thickness = 2
		height    = 7

		diameter = 56

		motorShaftDiameter = 5.2

		slots = 15
	)

	Register("encoder-disc", func() sdf.SDF3 {
		outerRing := NewCylinder(thickness, diameter+thickness*2).
			Difference(NewCylinder(thickness, diameter-thickness*2))

		encodeRing := NewCylinder(height, diameter+thickness).Difference(
			NewCylinder(height, diameter+thickness/2),
		).Difference(
			Slice(height, (diameter+thickness)/2, math.Pi/slots).RotateZCopy(slots),
		)

		arms := NewBox(diameter, thickness*2, thickness).RotateZCopy(4)
		shaft := NewCylinder(height-1.5, motorShaftDiameter+thickness*2)

		motorShaft := NewCylinder(height*2, motorShaftDiameter).Difference(
			NewBox(.5, motorShaftDiameter, height*2).Translate(motorShaftDiameter/2-.25, 0, 0),
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
