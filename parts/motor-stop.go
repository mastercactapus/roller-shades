package main

import (
	"github.com/deadsy/sdfx/sdf"
	builder "github.com/mastercactapus/sdf-builder"
)

func init() {
	const (
		w      = 10
		t      = 2
		oDia   = 22
		oT     = 2
		dia    = 6.2
		motorW = 42
	)
	Register("motor-stop", func() sdf.SDF3 {
		b := builder.NewBox(motorW+t*2, w, t+oT)
		b = b.Difference(
			builder.NewCylinder(b.SizeZ(), dia),
			builder.NewCylinder(oT, oDia).SnapMinZ(b.MinZ()),
		)

		leg := builder.NewBox(t, w, t).SnapMaxZ(b.MinZ())
		b = b.Union(
			leg.SnapMinX(b.MinX()),
			leg.SnapMaxX(b.MaxX()),
		)

		return b.Mirror(false, false, true).SnapMinZ(0)
	})
}
