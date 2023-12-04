package texture

import (
	"math"
	"ray_tracing/vector"
)

type Texture interface {
	Value(u, v float64, p vector.Point) vector.Color
}

type SolidColor struct {
	vector.Color
}

func NewSolidColor(c vector.Color) *SolidColor {
	return &SolidColor{
		Color: c,
	}
}

func (s *SolidColor) Value(u, v float64, p vector.Point) vector.Color {
	return s.Color
}

type CheckerTexture struct {
	invScale float64
	even     Texture
	odd      Texture
}

func NewCheckerTexture(scale float64, even, odd Texture) *CheckerTexture {
	return &CheckerTexture{
		invScale: 1.0 / scale,
		even:     even,
		odd:      odd,
	}
}

func (ct *CheckerTexture) Value(u, v float64, p vector.Point) vector.Color {

	xInt := int(math.Floor(ct.invScale * p.X()))
	yInt := int(math.Floor(ct.invScale * p.Y()))
	zInt := int(math.Floor(ct.invScale * p.Z()))
	isEven := (xInt+yInt+zInt)%2 == 0

	if isEven {
		return ct.even.Value(u, v, p)
	} else {
		return ct.odd.Value(u, v, p)
	}
}
