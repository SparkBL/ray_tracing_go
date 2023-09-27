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
