package main

import (
	"testing"
)

// benchmarking the decode function
func BenchmarkMH(b *testing.B) {
	for i := 0; i < b.N; i++ {
		monthHour()
	}
}

func BenchmarkGL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getLatest()
	}
}
