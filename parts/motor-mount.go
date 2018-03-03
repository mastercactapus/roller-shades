package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	"github.com/mastercactapus/sdf-builder"
)

func init() {
	const (
		t = 3
	)
	Register("motor-mount", func() sdf.SDF3 {
		h := cfg.Misc.Height
		screwD := cfg.Misc.ScrewDiameter
		piW := cfg.Inturrupter.Width
		piArmW := cfg.Inturrupter.ArmThickness
		piH := cfg.Inturrupter.BaseHeight
		piD := cfg.Inturrupter.Length
		piSlotW := cfg.Inturrupter.SlotWidth

		ledgeH := h - nema17Dia/2

		ledge := builder.NewBox(nema17Dia, t+1, ledgeH)

		motorScrewHole := builder.
			NewCylinder(t, screwD).
			RotateX(math.Pi / 2).
			SnapMidZ(ledge.MaxZ() + nema17ScrewOffset).
			SnapMinY(ledge.MaxY())

		mount := ledge.Union(
			builder.NewBox(nema17Dia, t, ledgeH+nema17ScrewOffset*2).SnapMinY(ledge.MaxY()),
		)
		piMount := builder.
			NewBox(piW, piH+1, piD).
			SnapMinY(0)
		piMount = piMount.
			Union(
				builder.NewBox(piW, t*2+1, piArmW*2+piSlotW).SnapMinY(0).SnapMaxZ(piMount.MaxZ()),
				builder.NewCylinder(t*2+1, 3).RotateX(math.Pi/2).SnapMidZ((piD-(piArmW*2+piSlotW))/2).SnapMinY(0),
			).
			Translate(0, 0, 1).
			SnapMinY(mount.MinY())

		mount = mount.Difference(
			piMount.RotateYOrigin(math.Pi/12, 0, 0, h),
			piMount.RotateYOrigin(-math.Pi/12, 0, 0, h),
		)

		pad := builder.NewBox(14, 14, 7).Difference(
			builder.
				NewCylinder(21, 9).
				SnapMaxZ(0).
				Union(
					builder.NewCone(5, 9, 4),
					builder.NewCylinder(20, 4),
				).
				Mirror(false, false, true).
				Translate(0, 0, 7),
		).SnapMaxY(mount.MaxY())

		mount = mount.
			Difference(
				builder.NewCylinder(nema17CenterRingH, nema17CenterRingDia).
					RotateX(math.Pi/2).
					SnapMidZ(h).
					SnapMinY(ledge.MaxY()),
				motorScrewHole.SnapMidX(ledge.MinX()+nema17ScrewOffset),
				motorScrewHole.SnapMidX(ledge.MaxX()-nema17ScrewOffset),
			).
			Union(
				builder.NewBox(t, 15, ledgeH).SnapMaxY(mount.MinY()),
				pad.SnapMaxX(mount.MinX()),
				pad.SnapMinX(mount.MaxX()),
			)

		return mount.RotateX(-math.Pi / 2).SnapMinZ(0)
	})
}
