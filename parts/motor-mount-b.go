package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	"github.com/mastercactapus/sdf-builder"
)

func init() {
	const (
		h = 40 // target height of the shaft center

		t      = 3
		screwD = 3.5

		// photointurrupter measurements. Expected = |Arm|Slot|Arm|Hole|
		piW          = 6.75 //width
		piD          = 15.6 // depth
		piH          = 2.5  // height of the base (not the arms)
		piSlotW      = 2    // width of the slot
		piArmW       = 4    // width of each arm
		piLeadOffset = 1    // offset from the edge that the leads stick out
		piHoleOffset = 2.25 // center of the mounting hole from the edge
	)
	Register("motor-mount-b", func() sdf.SDF3 {
		ledgeH := h - nema17Dia/2

		ledge := builder.NewBox(nema17Dia, t, ledgeH)

		motorScrewHole := builder.
			NewCylinder(t, screwD).
			RotateX(math.Pi / 2).
			SnapMidZ(ledge.MaxZ() + nema17ScrewOffset).
			SnapMinY(ledge.MaxY())

		mount := ledge.Union(
			builder.NewBox(nema17Dia, t, ledgeH+nema17ScrewOffset*2).SnapMinY(ledge.MaxY()),
			builder.NewBox(nema17Dia+t*5, t*2, t*2).SnapMinY(ledge.MinY()),
		)
		piMount := builder.
			NewBox(piW, piH, piD).
			Translate(0, 0, 1).
			SnapMinY(mount.MinY())

		mount = mount.Difference(
			piMount.RotateYOrigin(math.Pi/20, 0, 0, h),
			piMount.RotateYOrigin(-math.Pi/20, 0, 0, h),
		)

		screw := builder.
			NewCylinder(21, 9).
			SnapMaxZ(0).
			Union(
				builder.NewCone(5, 9, 4),
				builder.NewCylinder(20, 4),
			).
			Mirror(false, false, true).
			SnapMidY(mount.MidY()).
			Translate(0, 0, 9)

		screwOff := (mount.SizeY() - screw.SizeY()) / 2

		mount = mount.
			Difference(
				builder.NewCylinder(nema17CenterRingH, nema17CenterRingDia).
					RotateX(math.Pi/2).
					SnapMidZ(h).
					SnapMinY(ledge.MaxY()),
				motorScrewHole.SnapMidX(ledge.MinX()+nema17ScrewOffset),
				motorScrewHole.SnapMidX(ledge.MaxX()-nema17ScrewOffset),

				screw.SnapMinX(mount.MinX()+screwOff),
				screw.SnapMaxX(mount.MaxX()-screwOff),
			)

		return mount
		return mount.RotateX(-math.Pi / 2).SnapMinZ(0)
	})
}
