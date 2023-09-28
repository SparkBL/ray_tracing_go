package main

import (
	"math"
)

func DegressToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}
