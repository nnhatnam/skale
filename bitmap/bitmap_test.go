package bitmap

import (
	"math/rand"
	"testing"

	_ "net/http/pprof"
)

func BenchmarkSet(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(0))
	b.StartTimer()
	bm := NewBitMap(1000000)
	for i := 0; i < b.N; i++ {
		bm.Set(uint(r.Int31n(1000000)))
	}
}

func BenchmarkAll(b *testing.B) {

	bm := NewBitMap(1000000)
	b.StopTimer()
	r := rand.New(rand.NewSource(0))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bm.Set(uint(r.Int31n(1000000)))
	}
	bm.All()
}
