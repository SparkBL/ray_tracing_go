package interval

import "math"

type Interval [2]float64

func (i Interval) Min() float64 {
	return i[0]
}

func (i Interval) Max() float64 {
	return i[1]
}

func (i Interval) Size() float64 {
	return i[1] - i[0]
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

func (i Interval) Expand(delta float64) Interval {
	return Interval{
		i[0] - delta/2,
		i[1] + delta/2,
	}
}

func CombineIntervals(a, b Interval) Interval {
	return Interval{
		min(a[0], b[0]),
		max(a[1], b[1]),
	}
}

var Empty Interval = Interval{math.Inf(1), math.Inf(-1)}
var Universe Interval = Interval{math.Inf(-1), math.Inf(1)}
