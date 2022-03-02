package tests

import (
	"KeyValueHTTPStore/internal/cluster"
	"KeyValueHTTPStore/internal/command"
	"context"
	"runtime"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func Benchmark_Performance_Get(t *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	storage := cluster.New()
	storage.Initialize(ctx)

	cmds := []command.Command{
		command.Get("foo"),
	}

	for i := 0; i < t.N; i++ {
		storage.Execute(cmds)
	}

	cancel()
}

func Benchmark_Performance_Set(t *testing.B) {
	ctx, _ := context.WithCancel(context.Background())
	storage := cluster.New()
	storage.Initialize(ctx)

	cmds := []command.Command{
		command.Set("foo", "bar"),
	}

	for i := 0; i < t.N; i++ {
		storage.Execute(cmds)
	}
}
