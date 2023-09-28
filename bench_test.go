package main

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkSample(b *testing.B) {
	b.N = 10
	for i := 0; i < b.N; i++ {
		main()
	}
	fmt.Println(time.UnixMilli(b.Elapsed().Milliseconds() / int64(b.N)))
}
