package algorithms

import (
	"github.com/rafaelromao/labs/solutions/go/nearby"
	"github.com/rafaelromao/labs/solutions/go/ttlcache"
	"github.com/rafaelromao/labs/solutions/go/workerpool"
	"testing"
	"time"
)

var benchmarkResult any

func BenchmarkTopK(b *testing.B) {
	values := make([]int64, 10000)
	for i := range values {
		values[i] = int64((i * 7919) % 10007)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		benchmarkResult = TopK(values, 20)
	}
}
func BenchmarkSpatialRadius(b *testing.B) {
	idx := nearby.New()
	for i := 0; i < 10000; i++ {
		idx.Put(string(rune(i)), nearby.Point{Latitude: float64(i%180) - 90, Longitude: float64(i%360) - 180})
	}
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result, err := idx.Radius(nearby.Point{}, 10000)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkResult = result
	}
}
func BenchmarkTTLRead(b *testing.B) {
	cache, err := ttlcache.New[string, int](1000, time.Hour)
	if err != nil {
		b.Fatal(err)
	}
	cache.Put("key", 7)
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v, ok := cache.Get("key")
		if !ok {
			b.Fatal("miss")
		}
		benchmarkResult = v
	}
}
func BenchmarkWorkerSaturation(b *testing.B) {
	p, _ := workerpool.New(1, 1)
	started, release := make(chan struct{}), make(chan struct{})
	p.Submit(func() error { close(started); <-release; return nil })
	<-started
	p.Submit(func() error { return nil })
	defer func() { close(release); p.Close(); <-p.Done() }()
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, err := p.Submit(func() error { return nil }); err != workerpool.ErrFull {
			b.Fatal(err)
		}
	}
}
