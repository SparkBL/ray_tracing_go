package main

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

var Empty Interval = Interval{math.Inf(1), math.Inf(-1)}
var Universe Interval = Interval{math.Inf(-1), math.Inf(1)}
