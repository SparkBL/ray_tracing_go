package interval

import "math"

type Interval struct {
	Min, Max float64
}

func (i *Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}

func (i *Interval) Surrounds(x float64) bool {
	return i.Min < x && x < i.Max
}

func (i *Interval) Clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	}
	if x > i.Max {
		return i.Max
	}
	return x
}

var Empty Interval = Interval{math.Inf(1), math.Inf(-1)}
var Universe Interval = Interval{math.Inf(-1), math.Inf(1)}