package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

var cfg struct {
	Motor struct {
		ShaftDiameter float64
		FlatDepth     float64
	}
	Encoder struct {
		Diameter  float64
		Slots     int
		Thickness float64
		Height    float64
	}
	Misc struct {
		ScrewDiameter float64
		NutSize       float64
		Height        float64
	}
	Roller struct {
		ID float64
	}
	Inturrupter struct {
		Width        float64
		Length       float64
		BaseHeight   float64
		ArmThickness float64
		SlotWidth    float64
	}
	MagStop struct {
		MagnetDiameter float64
		RingDiameter   float64
		StopPositions  int

		Bearing struct {
			ID float64
			OD float64
			T  float64
		}
	}
}

func init() {
	_, err := toml.DecodeFile("config.toml", &cfg)
	if err != nil {
		log.Fatalln(err)
	}
	if cfg.MagStop.StopPositions%3 != 0 {
		log.Fatalln("MagStop.StopPositions must be divisible by 3.")
	}
}
