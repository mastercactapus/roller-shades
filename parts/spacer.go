package main

import (
	"github.com/deadsy/sdfx/sdf"
	"github.com/mastercactapus/sdf-builder"
)

func init() {
	const (
		inset = 1.5
		space = 1.5
	)
	Register("spacer", func() sdf.SDF3 {
		screwD := cfg.Misc.ScrewDiameter
		ID := cfg.MagStop.Bearing.ID
		OD := ID + 4
		return builder.
			NewCylinder(space+inset, ID).
			Union(
				builder.NewCylinder(space, OD),
			).
			Difference(
				builder.NewCylinder(space+inset, screwD),
			)
	})
}
