package ray

import "ray_tracing/vector"

type Ray struct {
	Origin    vector.Point
	Direction vector.Vector
}

func (r *Ray) At(t float64) vector.Point {
	return vector.Add(r.Origin, r.Direction.Multiply(t))
}
