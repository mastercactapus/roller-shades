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
		magRingCount = 10
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
			Translate(0, 0, 9)

		mount = mount.
			Union(
				builder.NewBox(w, w, h).SnapMaxX(mount.MaxX()),
				builder.NewBox(w, w, h).SnapMinX(mount.MinX()),
				builder.NewBox(dia+w*4+t*2, w, w),

				builder.NewCylinder(w, magRingDia+t*2).
					Difference(builder.NewCylinder(w, magRingDia-t*2)).
					RotateX(math.Pi/2).
					SnapMidZ(h).
					SnapMidY(0).
					Intersection(
						builder.NewBox(dia, w, h),
					),
			).
			Difference(
				builder.
					NewCylinder(w, dia).
					RotateX(math.Pi/2).
					SnapMidY(0).
					SnapMidZ(h),
			).
			Union(
				magMount.
					RotateYOrigin(math.Pi/magRingCount*2, 0, 0, h),
				magMount.
					RotateYOrigin(-math.Pi/magRingCount*2, 0, 0, h),
				magMount,

				builder.
					NewCylinder(1, dia+t*2).
					Difference(builder.NewCone(1, dia-2, dia)). // 1x1 = quarter angle or 45 deg. most printers should handle it fine
					RotateX(math.Pi/2).
					SnapMidZ(h).
					SnapMinY(mount.MaxY()),
			).
			Difference(
				magHole.
					RotateYOrigin(math.Pi/magRingCount*2, 0, 0, h),
				magHole.
					RotateYOrigin(-math.Pi/magRingCount*2, 0, 0, h),
				magHole,
			)

		mount = mount.Difference(
			screw.SnapMinX(mount.MinX()),
			screw.SnapMaxX(mount.MaxX()),
		)

		return mount.RotateX(math.Pi / 2).SnapMinZ(0)
	})
}
