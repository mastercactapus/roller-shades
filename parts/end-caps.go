package main

import (
	"github.com/mastercactapus/sdf-builder"

	"github.com/deadsy/sdfx/sdf"
)

func init() {

	const (
		t = 2
		w = 5
		h = 20

		shaftH = 5
	)

	cap := func() builder.Builder {
		ringDia := cfg.Roller.ID - t*2
		return builder.NewCylinder(shaftH, ringDia).
			Union(
				builder.NewCylinder(3, cfg.Roller.ID+t*2),
				builder.NewCylinder(h, ringDia).Difference(builder.NewCylinder(h, ringDia-t*2)),
				builder.
					NewRoundedBox(t*2, t*2, h+t*2, t).
					SnapMidX(ringDia/2).
					RotateZCopy(3).
					Translate(0, 0, -t).
					Difference(
						builder.NewCylinder(t, cfg.Roller.ID).SnapMaxZ(0),
					),
			)

	}
	Register("motor-cap", func() sdf.SDF3 {
		center := cap().
			Difference(
				builder.NewCylinder(shaftH, cfg.Motor.ShaftDiameter),
			).
			Union(
				builder.NewBox(cfg.Motor.FlatDepth, cfg.Motor.ShaftDiameter/1.5, shaftH).SnapMaxX(cfg.Motor.ShaftDiameter / 2),
			)

		return center
	})

	Register("idle-cap", func() sdf.SDF3 {
		center := cap().
			Difference(
				builder.NewCylinder(shaftH, cfg.Misc.ScrewDiameter),
				builder.NewHexagon(cfg.Misc.NutSize, 2.5).SnapMaxZ(shaftH),
			)

		return center
	})

}
