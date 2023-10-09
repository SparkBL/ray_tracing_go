package vector

import (
	"math"
	"math/rand"
	"sync"

	"github.com/seehuhn/mt19937"
)

func UnitVector(v Vector) Vector {
	return v.Divide(v.Length())

}

var r *rand.Rand = rand.New(mt19937.New())
var mu sync.Mutex

func Random() Vector {
	return Vector{
		r.Float64(),
		r.Float64(),
		r.Float64(),
	}
}

func RandomBounded(min, max float64) Vector {
	mu.Lock()
	defer mu.Unlock()
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

func RandomInUnitDisk() Vector {
	for {
		v := RandomBounded(-1, 1)
		v[2] = 0
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

func Refract(v, n Vector, et float64) Vector {
	cosTheta := math.Min(Dot(v.Negative(), n), 1.0)
	rOutPerp := v.Add(n.Multiply(cosTheta)).Multiply(et)
	rOutParallel := n.Multiply(-math.Sqrt(math.Abs(1 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}
