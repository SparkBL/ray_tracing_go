package concrand

import (
	"math/rand"

	"github.com/seehuhn/mt19937"
)

var r *ConcRand = NewConcaRand()

func Float64() float64 {
	return <-r.output
}

type ConcRand struct {
	*rand.Rand
	output chan float64
}

func NewConcaRand() *ConcRand {
	ret := &ConcRand{
		Rand:   rand.New(mt19937.New()),
		output: make(chan float64, 100000000),
	}
	for i := 0; i < cap(ret.output); i++ {
		ret.output <- ret.Rand.Float64()
	}

	go func() {
		for {
			//if len(ret.output) < cap(ret.output)/2 {
			//	for i := 0; i < (cap(ret.output) - len(ret.output)); i++ {
			ret.output <- ret.Rand.Float64()

			//	}
			//	runtime.Gosched()
			//	}
		}
	}()

	return ret
}

func GetChan(c *ConcRand) <-chan float64 {
	return c.output
}

func Get(c *ConcRand) float64 {
	return <-c.output
}
