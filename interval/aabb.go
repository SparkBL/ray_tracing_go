package interval

import (
	"ray_tracing/ray"
	"ray_tracing/vector"
)

type AABB [3]Interval

type AABBOPtion func() AABB

func FromIntervals(x, y, z Interval) AABBOPtion {
	return func() AABB {
		return AABB{
			x,
			y,
			z,
		}
	}
}

func FromPoints(a, b vector.Point) AABBOPtion {
	return func() AABB {
		return AABB{
			Interval{min(a[0], b[0]), max(a[0], b[0])},
			Interval{min(a[1], b[1]), max(a[1], b[1])},
			Interval{min(a[2], b[2]), max(a[2], b[2])},
		}
	}
}

func NewAABB(opt AABBOPtion) AABB {
	return opt()
}

func CombineAABB(a, b AABB) AABB {
	return AABB{
		CombineIntervals(a[0], b[0]),
		CombineIntervals(a[1], b[1]),
		CombineIntervals(a[2], b[2]),
	}
}

func (a *AABB) Axis(n int) Interval {
	if n > 2 {
		return a[0]
	}
	return a[n]
}

// func (a *AABB) Hit(rIn *ray.Ray, rayT Interval) bool {
// 	for i, e := range a {
// 		t0 := min(
// 			e[0]-rIn.Origin[i]/rIn.Direction[i],
// 			e[1]-rIn.Origin[i]/rIn.Direction[i],
// 		)
// 		t1 := max(
// 			e[0]-rIn.Origin[i]/rIn.Direction[i],
// 			e[1]-rIn.Origin[i]/rIn.Direction[i],
// 		)
// 		rayT[0] = max(t0, rayT[0])
// 		rayT[1] = max(t1, rayT[1])
// 		if rayT[1] <= rayT[0] {
// 			return false
// 		}
// 	}
// 	return true
// }

func (a *AABB) Hit(rIn *ray.Ray, rayT Interval) bool {
	for i, e := range a {

		invD := 1 / rIn.Direction[i]
		orig := rIn.Origin[i]
		t0 := (e[0] - orig) * invD
		t1 := (e[1] - orig) * invD

		if t0 > rayT[0] {
			rayT[0] = t0
		}
		if t1 < rayT[1] {
			rayT[1] = t1
		}
		if rayT[1] <= rayT[0] {
			return false
		}
	}
	return true
}
