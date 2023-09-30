package util

import (
	"math"
)

func DegressToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func LinearToGamma(linearComponent float64) float64 {
	return math.Sqrt(linearComponent)
}
