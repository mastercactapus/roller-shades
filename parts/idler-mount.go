package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	"github.com/mastercactapus/sdf-builder"
)

func init() {
	const (
		t   = 2
		w   = 7
		dia = 22
		h   = 40

		magD         = 6.55
		magRingDia   = 40
		magRingCount = 9
	)
	Register("idler-mount", func() sdf.SDF3 {
		mount := builder.
			NewCylinder(w, dia+t*2).
			RotateX(math.Pi / 2).
			SnapMidY(0).
			SnapMidZ(h)

		magMount := builder.
			NewCylinder(w, magD+t*2).
			RotateX(math.Pi / 2).
			SnapMidY(0).
			SnapMidZ(h - magRingDia/2)

		magHole := builder.
			NewCylinder(w-.65, magD).
			RotateX(math.Pi / 2).
			SnapMaxY(mount.MaxY()).
			SnapMidZ(h - magRingDia/2)

		screw := builder.
			NewCylinder(21, 9).
			SnapMaxZ(0).
			Union(
				builder.NewCone(5, 9, 4),
				builder.NewCylinder(20, 4),
			).
			Mirror(false, false, true).
			SnapMidY(mount.MidY()).
			Translate(0, 0, 8)

		foot := builder.
			NewBox(dia+w*4+t*2, w, w)

		toe := builder.
			NewBox(w*2, w*2, w).
			Difference(screw).
			SnapMinY(foot.MaxY())
		foot = foot.
			Union(
				toe.SnapMinX(foot.MinX()),
				toe.SnapMaxX(foot.MaxX()),
			)

		mount = mount.
			Union(
				builder.NewBox(w, w, h).SnapMaxX(mount.MaxX()),
				builder.NewBox(w, w, h).SnapMinX(mount.MinX()),
				foot,

				builder.NewCylinder(w, magRingDia+t*2).
					Difference(builder.NewCylinder(w, magRingDia-t*2)).
					RotateX(math.Pi/2).
					SnapMidZ(h).
					SnapMidY(0),

				builder.NewBox(dia+t*2-w*2, w, (magRingDia-dia)/2-t).SnapMinZ(mount.MaxZ()-t),
			).
			Difference(
				builder.
					NewCylinder(w, dia).
					RotateX(math.Pi/2).
					SnapMidY(0).
					SnapMidZ(h),
			).
			Union(
				magMount.RotateYCopyOrigin(3, 0, 0, h),
				builder.
					NewCylinder(1, dia+t*2).
					Difference(builder.NewCone(1, dia-2, dia)). // 1x1 = quarter angle or 45 deg. most printers should handle it fine
					RotateX(math.Pi/2).
					SnapMidZ(h).
					SnapMinY(mount.MaxY()),
			).
			Difference(
				magHole.RotateYCopyOrigin(3, 0, 0, h),
			)

		return mount.RotateX(math.Pi / 2).SnapMinZ(0)
	})
}
