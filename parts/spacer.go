package main

import (
	"github.com/deadsy/sdfx/sdf"
	"github.com/mastercactapus/sdf-builder"
)

func init() {
	const (
		OD     = 12
		ID     = 8
		screwD = 3.5
		inset  = 1.5
		space  = .75
	)
	Register("spacer", func() sdf.SDF3 {
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
