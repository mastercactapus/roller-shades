package main

import (
	"github.com/deadsy/sdfx/sdf"
	"github.com/mastercactapus/sdf-builder"
)

func init() {
	const (
		magH = 1.6
		magD = 6.45

		n = 9

		dia     = 40 // diameter of magnetic ring
		padding = 2
		magT    = 0.65

		screwD = 3.5
		h      = 3
	)

	Register("magnetic-stop", func() sdf.SDF3 {
		ring := builder.NewCylinder(h, dia+magD+padding*2)
		ring = ring.Difference(
			builder.NewCylinder(h, screwD),

			builder.
				NewCylinder(h-magT, magD).
				Translate(dia/2, 0, 0).
				SnapMaxZ(ring.MaxZ()).
				RotateZCopy(n),
		)
		return ring
	})

}
