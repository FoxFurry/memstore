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

func Benchmark_Performance(t *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	storage := cluster.New()
	storage.Initialize(ctx)

	cmds := []command.Command{
		command.Set("test", "testv"),
	}

	for i := 0; i < t.N; i++ {
		storage.Execute(cmds)
	}

	cancel()
}
