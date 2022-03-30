package tests

import (
	"context"
	"fmt"
	"github.com/FoxFurry/GoKeyValueStore/internal/cluster"
	"github.com/FoxFurry/GoKeyValueStore/internal/command"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	keyMaxLen   = 20
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())
}

func randString() string {
	n := rand.Intn(keyMaxLen)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

var testSizes = []int{
	1,
	10,
	100,
	1000,
	10000,
}

func Benchmark_Get_SameData(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	storage := cluster.New()
	storage.Initialize(ctx)

	for _, size := range testSizes {
		b.Run(fmt.Sprintf("%s-%d", b.Name(), size), func(b *testing.B) {
			var cmds []command.Command

			for i := 0; i < size; i++ {
				cmds = append(cmds, command.Get("foo"))
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				storage.Execute(cmds)
			}
		})
	}

	cancel()
}

func Benchmark_Set_SameData(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	storage := cluster.New()
	storage.Initialize(ctx)

	for _, size := range testSizes {
		b.Run(fmt.Sprintf("%s-%d", b.Name(), size), func(b *testing.B) {
			var cmds []command.Command

			for i := 0; i < size; i++ {
				cmds = append(cmds, command.Set("foo", "bar"))
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				storage.Execute(cmds)
			}
		})
	}

	cancel()
}

func Benchmark_Get_VariousData(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	storage := cluster.New()
	storage.Initialize(ctx)

	for _, size := range testSizes {
		b.Run(fmt.Sprintf("%s-%d", b.Name(), size), func(b *testing.B) {
			var cmds []command.Command

			for i := 0; i < size; i++ {
				cmds = append(cmds, command.Get(randString()))
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				storage.Execute(cmds)
			}
		})
	}

	cancel()
}

func Benchmark_Set_VariousData(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	storage := cluster.New()
	storage.Initialize(ctx)

	for _, size := range testSizes {
		b.Run(fmt.Sprintf("%s-%d", b.Name(), size), func(b *testing.B) {
			var cmds []command.Command

			for i := 0; i < size; i++ {
				cmds = append(cmds, command.Set(randString(), randString()))
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				storage.Execute(cmds)
			}
		})
	}

	cancel()
}
