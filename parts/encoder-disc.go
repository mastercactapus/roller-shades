package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	"github.com/mastercactapus/sdf-builder"
)

func init() {
	Register("encoder-disc", func() sdf.SDF3 {
		outerRing := builder.
			NewCylinder(cfg.Encoder.Thickness, cfg.Encoder.Diameter+cfg.Encoder.Thickness*2).
			Difference(
				builder.NewCylinder(cfg.Encoder.Thickness, cfg.Encoder.Diameter-cfg.Encoder.Thickness*2),
			)

		encodeRing := builder.NewCylinder(cfg.Encoder.Height, cfg.Encoder.Diameter+cfg.Encoder.Thickness).Difference(
			builder.NewCylinder(cfg.Encoder.Height, cfg.Encoder.Diameter+cfg.Encoder.Thickness/2),
		).Difference(
			builder.NewSlice(cfg.Encoder.Height, (cfg.Encoder.Diameter+cfg.Encoder.Thickness)/2, math.Pi/float64(cfg.Encoder.Slots)).
				RotateZCopy(cfg.Encoder.Slots),
		)

		arms := builder.NewBox(cfg.Encoder.Diameter, cfg.Encoder.Thickness*2, cfg.Encoder.Thickness).
			RotateZCopy(4)
		shaft := builder.NewCylinder(cfg.Encoder.Height-1.5, cfg.Motor.ShaftDiameter+cfg.Encoder.Thickness*2)

		motorShaft := builder.
			NewCylinder(cfg.Encoder.Height*2, cfg.Motor.ShaftDiameter).
			Difference(
				builder.NewBox(cfg.Motor.FlatDepth, cfg.Motor.ShaftDiameter, cfg.Encoder.Height*2).
					Translate(cfg.Motor.ShaftDiameter/2-cfg.Motor.FlatDepth/2, 0, 0),
			)

		mount := arms.
			Union(shaft).
			Difference(motorShaft)
		f := sdf.Union3D(
			outerRing,
			encodeRing,
			mount,
		)
		return f
	})
}
