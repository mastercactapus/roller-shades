package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	builder "github.com/mastercactapus/sdf-builder"
)

func init() {
	const (
		sensW = 9.0
		sensH = 5.5
		sensT = 2.0

		t = 10

		motorD = 37
		screwD = 31.5
		depth  = 40
		h      = 40
	)
	Register("motor-mount-a", func() sdf.SDF3 {
		sh := sensH + 6
		sw := sensW + 3
		intMount := builder.NewBox(sw, sensT, sh)
		intMount = intMount.Difference(
			builder.NewBox(sensW, sensT+2, sensH).SnapMidZ(intMount.MidZ()),
		).SnapMinX(0)

		intArm := builder.NewRTriangle(sw, t-sensT, 2).SnapMinX(0).SnapMinY(intMount.MaxY())
		intArm = intArm.Union(
			builder.NewBox(1, t-sensT, 2).SnapMaxX(0).SnapMaxY(intArm.MaxY()),
		)

		intMount = intMount.Union(
			intArm,
			intArm.SnapMaxZ(intMount.MaxZ()),
			builder.NewBox(1, sensT, sh).SnapMinX(-1),
		)

		post := builder.NewBox(t, t, h+t/2+1).
			Difference(
				builder.NewCylinder(t, 3.5).SnapMidZ(0).
					RotateX(-math.Pi/2).
					Translate(0, 0, h),
			).Translate(screwD/2, 0, 0)

		foot := builder.NewBox(screwD/2+t*2, t, t).SnapMinX(0)
		screw := builder.NewCylinder(21, 9).SnapMaxZ(0).Union(
			builder.NewCone(5, 9, 4),
			builder.NewCylinder(20, 4),
		).Mirror(false, false, true)
		foot = foot.Difference(
			screw.
				SnapMinX(foot.MaxX()-screw.SizeX()-(foot.SizeY()-screw.SizeY())/2).
				Translate(0, 0, t),
		)
		post = post.Union(
			intMount.SnapMinX(post.MaxX()-1).SnapMaxY(post.MaxY()).SnapMidZ(h).RotateYOrigin(math.Pi/41, 0, 0, h),
			foot,
		).Difference(
			builder.NewCylinder(t, motorD).
				RotateX(-math.Pi / 2).
				SnapMidZ(h),
		)

		return post.Union(
			post.Mirror(true, false, false),
		).RotateX(math.Pi / 2).SnapMinZ(0).SnapMidY(0)
	})
}
