package examples

import (
	"github.com/scarabsoft/go-snowflake"
	"runtime"
	"sync"
	"testing"
)

var cores = runtime.NumCPU()
var wg = sync.WaitGroup{}
var gen, _ = snowflake.NewGenerator()

func BenchmarkTestBenchmark_Single(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = gen.Next()
	}
}

func BenchmarkTestBenchmark_Parallel(b *testing.B) {
	for i := 0; i < cores; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N/cores; i++ {
				_, _ = gen.Next()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
