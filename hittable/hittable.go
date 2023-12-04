package hittable

import (
	"math"
	"ray_tracing/interval"
	"ray_tracing/ray"
	"ray_tracing/vector"
	"sort"
)

type HitRecord struct {
	Point       vector.Point
	Normal      vector.Vector
	Material    Material
	T           float64
	IsFrontFace bool
}

func (hr *HitRecord) SetFaceNormal(r *ray.Ray, outwardNormal vector.Vector) {
	// Sets the hit record normal vector.
	// NOTE: the parameter `outward_normal` is assumed to have unit length.

	hr.IsFrontFace = vector.Dot(r.Direction, outwardNormal) < 0
	if hr.IsFrontFace {
		hr.Normal = outwardNormal
	} else {
		hr.Normal = outwardNormal.Negative()
	}

}

type Hittable interface {
	Hit(r *ray.Ray, rayT interval.Interval, rec *HitRecord) bool
	BoundingBox() interval.AABB
}

type Sphere struct {
	Center          vector.Point
	Radius          float64
	Material        Material
	move            bool
	targetDirection vector.Vector
	bbox            interval.AABB
}

func NewSphere(center vector.Point, radius float64, material Material) *Sphere {
	rvec := vector.Vector{radius, radius, radius}
	return &Sphere{
		Center:          center,
		Radius:          radius,
		Material:        material,
		move:            false,
		targetDirection: center,
		bbox: interval.NewAABB(interval.FromPoints(
			center.Add(rvec.Negative()),
			center.Add(rvec),
		)),
	}
}

func (s *Sphere) BoundingBox() interval.AABB {
	return s.bbox
}

func (s *Sphere) MoveTo(to vector.Point) {
	s.move = true
	s.targetDirection = to.Add(s.Center.Negative())
	rvec := vector.Vector{s.Radius, s.Radius, s.Radius}
	s.bbox = interval.CombineAABB(
		s.bbox,
		interval.NewAABB(interval.FromPoints(
			to.Add(rvec.Negative()),
			to.Add(rvec),
		)))

}

func (s *Sphere) CenterAt(time float64) vector.Point {
	if !s.move {
		return s.Center
	}
	return s.Center.Add(s.targetDirection.Multiply(time))
}

func (s *Sphere) Hit(r *ray.Ray, rayT interval.Interval, rec *HitRecord) bool {
	center := s.CenterAt(r.Time)
	ocDistance := r.Origin.Add(center.Negative())
	a := r.Direction.LengthSquared()
	halfB := vector.Dot(ocDistance, r.Direction)
	c := ocDistance.LengthSquared() - s.Radius*s.Radius
	discriminant := halfB*halfB - a*c
	if discriminant < 0 {
		return false
	}
	// Find the nearest root that lies in the acceptable range.
	sqrtd := math.Sqrt(discriminant)
	root := (-halfB - sqrtd) / a

	if !rayT.Surrounds(root) {
		root = (-halfB + sqrtd) / a
		if !rayT.Surrounds(root) {
			return false
		}
	}
	rec.T = root
	rec.Point = r.At(rec.T)
	outwardNormal := rec.Point.Add(s.Center.Negative()).Divide(s.Radius)
	rec.SetFaceNormal(r, outwardNormal)
	rec.Material = s.Material
	return true
}

type Plane struct {
	Center   vector.Point
	Normal   vector.Vector
	Material Material
}

func (s *Plane) Hit(r *ray.Ray, rayT interval.Interval, rec *HitRecord) bool {
	denominator := vector.Dot(s.Normal, r.Direction)
	if math.Abs(denominator) < 0.0 {
		return false
	}
	root := vector.Dot(s.Center.Add(r.Origin.Negative()), s.Normal) / denominator
	if !rayT.Surrounds(root) {
		return false
	}
	rec.T = root
	rec.Point = r.At(rec.T)
	outwardNormal := rec.Point.Add(s.Center.Negative()).Divide(denominator)
	rec.SetFaceNormal(r, outwardNormal)
	rec.Material = s.Material

	return true
}

type Hittables struct {
	objects []Hittable
	bbox    interval.AABB
}

func (hl *Hittables) BoundingBox() interval.AABB {
	return hl.bbox
}

func (hl *Hittables) Append(objects ...Hittable) {
	hl.objects = append(hl.objects, objects...)
	for _, o := range objects {
		hl.bbox = interval.CombineAABB(hl.bbox, o.BoundingBox())
	}
}

func (hl *Hittables) ToBVHTree() *BVHNode {
	return NewBHVTree(hl.objects...)
}

func NewWorld(o ...Hittable) *Hittables {
	return &Hittables{
		objects: o,
	}
}

func (hl *Hittables) Hit(r *ray.Ray, rayT interval.Interval, rec *HitRecord) bool {
	hitAnything := false
	closestSoFar := rayT.Max()
	for _, h := range hl.objects {
		if h.Hit(r, interval.Interval{rayT.Min(), closestSoFar}, rec) {
			hitAnything = true
			closestSoFar = rec.T
		}
	}
	return hitAnything
}

func boxCompare(h1, h2 Hittable, axis int) bool {
	return h1.BoundingBox().Axis(axis).Min() < h2.BoundingBox().Axis(axis).Min()
}

type BVHNode struct {
	left, right Hittable
	bbox        interval.AABB
}

func NewBHVTree(src ...Hittable) *BVHNode {
	axis := randGen.Intn(2)

	node := BVHNode{}
	switch len(src) {
	case 1:
		node.left, node.right = src[0], src[0]
	case 2:
		if boxCompare(src[0], src[1], axis) {
			node.left = src[0]
			node.right = src[1]
		} else {
			node.left = src[1]
			node.right = src[0]
		}
	default:
		sort.Slice(src, func(i, j int) bool { return boxCompare(src[i], src[j], axis) })
		middle := len(src) / 2
		node.left = NewBHVTree(src[:middle]...)
		node.right = NewBHVTree(src[middle:]...)
	}
	node.bbox = interval.CombineAABB(node.left.BoundingBox(), node.right.BoundingBox())
	return &node
}

func (b *BVHNode) BoundingBox() interval.AABB {
	return b.bbox
}

func (b *BVHNode) Hit(r *ray.Ray, rayT interval.Interval, rec *HitRecord) bool {
	if !b.bbox.Hit(r, rayT) {
		return false
	}

	var maxT float64
	hitLeft := b.left.Hit(r, rayT, rec)
	if hitLeft {
		maxT = rec.T
	} else {
		maxT = rayT.Max()
	}
	hitRight := b.right.Hit(r, interval.Interval{rayT.Min(), maxT}, rec)

	return hitLeft || hitRight
}
