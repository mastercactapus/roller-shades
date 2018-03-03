package main

import (
	"github.com/deadsy/sdfx/sdf"
	"github.com/mastercactapus/sdf-builder"
)

func init() {
	const (
		padding = 2
		magT    = 0.65

		h = 3
	)

	Register("magnetic-stop", func() sdf.SDF3 {
		ring := builder.NewCylinder(h, cfg.MagStop.RingDiameter+cfg.MagStop.MagnetDiameter+padding*2)
		ring = ring.Difference(
			builder.NewCylinder(h, cfg.Misc.ScrewDiameter),

			builder.
				NewCylinder(h-magT, cfg.MagStop.MagnetDiameter).
				Translate(cfg.MagStop.RingDiameter/2, 0, 0).
				SnapMaxZ(ring.MaxZ()).
				RotateZCopy(cfg.MagStop.StopPositions),
		)
		return ring
	})

}
