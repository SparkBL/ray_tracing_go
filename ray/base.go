package ray

import (
	"ray_tracing/vector"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} { return &Ray{} },
}

func getRay() (r *Ray) {
	ifc := pool.Get()
	if ifc != nil {
		r = ifc.(*Ray)
	}
	return
}

func putRay(r *Ray) {
	r.Origin = vector.Vector{0, 0, 0}
	r.Direction = vector.Vector{0, 0, 0}
	pool.Put(r)
}

type Ray struct {
	Origin    vector.Point
	Direction vector.Vector
}

func New() *Ray {
	return getRay()
}

func Save(r *Ray) {
	putRay(r)
}

func (r *Ray) At(t float64) vector.Point {
	return r.Origin.Add(r.Direction.Multiply(t))
}
