package ray

import "ray_tracing/vector"

type Ray struct {
	Origin    vector.Point
	Direction vector.Vector
}

func (r *Ray) At(t float64) vector.Point {
	return r.Origin.Add(r.Direction.Multiply(t))
}
