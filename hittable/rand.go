package hittable

import (
	"math/rand"

	"github.com/seehuhn/mt19937"
)

var randGenerator *rand.Rand = rand.New(mt19937.New())
