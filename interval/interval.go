package interval

import "math"

type Interval [2]float64

func (i *Interval) Min() float64 {
	return i[0]
}

func (i *Interval) Max() float64 {
	return i[1]
}

func (i *Interval) Contains(x float64) bool {
	return i[0] <= x && x <= i[1]
}

func (i *Interval) Surrounds(x float64) bool {
	return i[0] < x && x < i[1]
}

func (i *Interval) Clamp(x float64) float64 {
	if x < i[0] {
		return i[0]
	}
	if x > i[1] {
		return i[1]
	}
	return x
}

var Empty Interval = Interval{math.Inf(1), math.Inf(-1)}
var Universe Interval = Interval{math.Inf(-1), math.Inf(1)}
