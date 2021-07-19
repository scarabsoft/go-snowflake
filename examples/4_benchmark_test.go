package examples

import (
	"github.com/scarabsoft/go-snowflake"
	"runtime"
	"sync"
	"testing"
)

func BenchmarkTestBenchmark_Single(b *testing.B) {
	gen, _ := snowflake.New()

	for i := 0; i < b.N; i++ {
		_ = <- gen.Next()
	}
}

func BenchmarkTestBenchmark_Multiple(b *testing.B) {
	cores := runtime.NumCPU()

	gen, _ := snowflake.New()

	wg := sync.WaitGroup{}
	for i := 0; i < cores; i++ {
		go func() {
			wg.Add(1)
			for i := 0; i < b.N/cores; i++ {
				gen.Next()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
