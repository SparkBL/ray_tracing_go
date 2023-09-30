package vector

import (
	"math/rand"

	"github.com/seehuhn/mt19937"
)

func UnitVector(v Vector) Vector {
	return v.Divide(v.Length())

}

var r *rand.Rand = rand.New(mt19937.New())

func Random() Vector {
	return Vector{
		r.Float64(),
		r.Float64(),
		r.Float64(),
	}
}

func RandomBounded(min, max float64) Vector {
	return Vector{
		min + r.Float64()*(max-min),
		min + r.Float64()*(max-min),
		min + r.Float64()*(max-min),
	}
}

func RandomInUnitSphere() Vector {
	for {
		v := RandomBounded(-1, 1)
		if v.LengthSquared() < 1 {
			return v
		}
	}
}

func RandomUnitVector() Vector {
	return UnitVector(RandomInUnitSphere())
}

func RandomOnHemisphere(normal Vector) Vector {
	ushp := RandomUnitVector()
	if Dot(ushp, normal) > 0.0 {
		return ushp
	} else {
		return ushp.Negative()
	}
}

func Reflect(v, n Vector) Vector {
	return v.Add(n.Multiply(Dot(v, n) * 2).Negative())
}
