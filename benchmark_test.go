package log8q

import (
	"context"
	"testing"
)

func Benchmark_test(B *testing.B) {
	B.ReportAllocs()
	B.ResetTimer()
	log := New(context.Background(), nil)
	for i := 0; i < B.N; i++ {
		log.Info("bench testing data for every day")
		log.Infof("bench testing data for every day! %s", "OK-opsdfds")
	}
}
