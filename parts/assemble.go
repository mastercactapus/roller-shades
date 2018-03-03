package main

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	"github.com/mastercactapus/sdf-builder"
)

func assembleMesh() sdf.SDF3 {
	const (
		motorLen  = 35
		shaftLen  = 25
		grooveLen = 18
	)

	m3 := func(size float64) builder.Builder {
		return builder.
			NewCylinder(size, 3).
			Union(
				builder.
					NewCylinder(3, 5.6).
					SnapMaxZ(0),
			)
	}
	m3Washer := func(size float64) builder.Builder {
		return m3(size).Union(
			builder.
				NewCylinder(.75, 7).
				Difference(
					builder.NewCylinder(.75, cfg.Misc.ScrewDiameter),
				).
				Translate(0, 0, size+1),
		)
	}

	nut := builder.
		NewHexagon(cfg.Misc.NutSize, 3).
		Difference(
			builder.NewCylinder(3, cfg.Misc.ScrewDiameter),
		).
		RotateY(math.Pi / 2)

	motor := builder.
		NewBox(nema17Dia, nema17Dia, motorLen)

	screw := builder.
		NewCylinder(5, cfg.Misc.ScrewDiameter).
		SnapMaxZ(motor.MaxZ()).
		SnapMidY(motor.MinY() + nema17ScrewOffset)

	motor = motor.
		Difference(
			screw.SnapMidX(motor.MinX()+nema17ScrewOffset),
			screw.SnapMidX(motor.MaxX()-nema17ScrewOffset),
			screw.SnapMidX(motor.MinX()+nema17ScrewOffset).SnapMidY(motor.MaxY()-nema17ScrewOffset),
			screw.SnapMidX(motor.MaxX()-nema17ScrewOffset).SnapMidY(motor.MaxY()-nema17ScrewOffset),
		).
		Union(
			builder.
				NewCylinder(nema17CenterRingH, nema17CenterRingDia).
				SnapMinZ(motor.MaxZ()),
			builder.
				NewCylinder(shaftLen, nema17ShaftDia).
				SnapMinZ(motor.MaxZ()),
		).
		RotateY(math.Pi / 2)

	get := func(name string) builder.Builder {
		return builder.Builder{SDF3: getMesh(name)}
	}

	mount := get("motor-mount").
		RotateX(math.Pi / 2).
		RotateZ(-math.Pi / 2)

	pi := builder.
		NewBox(cfg.Inturrupter.BaseHeight, cfg.Inturrupter.Width, cfg.Inturrupter.Length)
	piArm := builder.
		NewBox(6, cfg.Inturrupter.Width, cfg.Inturrupter.ArmThickness).
		SnapMinX(pi.MaxX())
	pi = pi.
		Union(
			piArm.SnapMaxZ(pi.MaxZ()),
			piArm.SnapMaxZ(pi.MaxZ()-piArm.SizeZ()-cfg.Inturrupter.SlotWidth),
			m3(8).RotateY(math.Pi/2).SnapMidZ(2.5).SnapMaxX(pi.MinX()-2),
			nut.SnapMidZ(2.5).SnapMinX(pi.MaxX()+18),
		).
		Difference(
			builder.
				NewCylinder(cfg.Inturrupter.BaseHeight, cfg.Misc.ScrewDiameter).
				RotateY(math.Pi/2).
				SnapMidZ(2.5).
				SnapMinX(pi.MinX()),
		).
		Translate(-17, 0, 0)

	enc := get("encoder-disc").RotateY(-math.Pi / 2)
	mCap := get("motor-cap").RotateY(math.Pi / 2)

	iCap := get("idle-cap").RotateY(-math.Pi / 2)
	iCap = iCap.
		Union(nut.SnapMaxX(iCap.MinX()))

	mStop := get("magnetic-stop")
	mStop = mStop.
		Union(
			builder.
				NewCylinder(1.5, cfg.MagStop.MagnetDiameter).
				Translate(cfg.MagStop.RingDiameter/2, 0, 0).
				TranslateCopy(0, 0, 2.5).
				RotateZCopy(cfg.MagStop.StopPositions).
				SnapMinZ(mStop.MaxZ() + 2),
		).
		RotateY(-math.Pi / 2)

	ringSpacer := get("spacer").RotateY(math.Pi / 2)

	bearing := builder.
		NewCylinder(cfg.MagStop.Bearing.T, cfg.MagStop.Bearing.OD).
		Difference(
			builder.NewCylinder(cfg.MagStop.Bearing.T, cfg.MagStop.Bearing.ID),
		).
		RotateY(math.Pi / 2)

	iMount := get("idler-mount")
	iMount = iMount.
		Union(
			builder.
				NewCylinder(1.5, cfg.MagStop.MagnetDiameter).
				Translate(0, -cfg.Misc.Height+cfg.MagStop.RingDiameter/2, 0).
				TranslateCopyN(3, 0, 0, 2.5).
				RotateZCopyOrigin(3, 0, -cfg.Misc.Height, 0).
				SnapMinZ(iMount.MaxZ() - 10),
		)

	iMount = iMount.
		RotateY(math.Pi / 2).
		RotateX(-math.Pi / 2)

	bushing := ringSpacer.
		Mirror(true, false, false)
	bushing = bushing.
		Union(
			m3Washer(20).RotateY(-math.Pi / 2).SnapMinX(bushing.MaxX() + 3),
		)

	// position things
	motor = motor.SnapMidZ(cfg.Misc.Height).SnapMaxX(mount.MinX() + shaftLen)
	enc = enc.SnapMidZ(cfg.Misc.Height).SnapMinX(motor.MaxX())
	mCap = mCap.SnapMidZ(cfg.Misc.Height).SnapMinX(enc.MaxX())
	iCap = iCap.SnapMidZ(cfg.Misc.Height).SnapMinX(mCap.MaxX() + 10)
	mStop = mStop.SnapMidZ(cfg.Misc.Height).SnapMinX(iCap.MaxX() - 3)
	ringSpacer = ringSpacer.SnapMidZ(cfg.Misc.Height).SnapMinX(mStop.MaxX() + 3)
	bearing = bearing.SnapMidZ(cfg.Misc.Height).SnapMinX(ringSpacer.MaxX() + 3)
	iMount = iMount.SnapMinX(bearing.MaxX() + 3)
	bushing = bushing.SnapMidZ(cfg.Misc.Height).SnapMinX(iMount.MaxX() - 10)

	mount = mount.Union(
		pi.RotateXOrigin(math.Pi/12, 0, 0, cfg.Misc.Height),
		pi.RotateXOrigin(-math.Pi/12, 0, 0, cfg.Misc.Height),
		m3Washer(6).
			RotateY(-math.Pi/2).
			SnapMidZ(motor.MinZ()+nema17ScrewOffset).
			SnapMidY(motor.MinY()+nema17ScrewOffset).
			SnapMinX(mount.MaxX()+3).
			TranslateCopy(0, nema17ScrewDist, 0),
	)

	return sdf.Union3D(
		mount,
		motor,
		enc,
		mCap,
		iCap,
		mStop,
		ringSpacer,
		bearing,
		iMount,
		bushing,
	)
}

// const (
// 	nema17Dia           = 42.3
// 	nema17ScrewDist     = 31
// 	nema17ScrewOffset   = (nema17Dia - nema17ScrewDist) / 2
// 	nema17ShaftDia      = 5
// 	nema17CenterRingDia = 22
// 	nema17CenterRingH   = 2
// )
