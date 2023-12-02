package hittable

import (
	"time"

	"golang.org/x/exp/rand"

	prng "gonum.org/v1/gonum/mathext/prng"
)

var randGen *rand.Rand = rand.New(prng.NewSplitMix64(uint64(time.Now().UnixNano())))
