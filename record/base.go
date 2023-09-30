package record

// type HitRecord struct {
// 	Point       vector.Point
// 	Normal      vector.Vector
// 	Material    material.Material
// 	T           float64
// 	IsFrontFace bool
// }

// func (hr *HitRecord) SetFaceNormal(r ray.Ray, outwardNormal vector.Vector) {
// 	// Sets the hit record normal vector.
// 	// NOTE: the parameter `outward_normal` is assumed to have unit length.

// 	hr.IsFrontFace = vector.Dot(r.Direction, outwardNormal) < 0
// 	if hr.IsFrontFace {
// 		hr.Normal = outwardNormal
// 	} else {
// 		hr.Normal = outwardNormal.Negative()
// 	}

// }
