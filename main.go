package main

import (
	"KeyValueHTTPStore/internal/cluster"
	"KeyValueHTTPStore/internal/command"
	"context"
	"fmt"
	"runtime"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	// Init
	ctx, cancel := context.WithCancel(context.Background())
	storage := cluster.New()
	storage.Initialize(ctx)
	time.Sleep(time.Nanosecond * 10000) // Let all goroutines start and get ready

	const N = 10
	var totalTime time.Duration
	for i := 0; i < N; i++ {

		// Execute set
		cmds := []command.Command{
			command.Set("foo", "bar"),
		}
		result, _ := storage.Execute(cmds)

		// Execute get
		cmds = []command.Command{
			command.Get("foo"),
		}

		setExec := time.Now() // Record time
		for {
			result, _ = storage.Execute(cmds)

			if result != nil && result[0] == "bar" {
				break
			}
		}

		elapsed := time.Since(setExec)
		totalTime += elapsed

		fmt.Printf("Try #%d: %s\n", i, elapsed)
	}

	fmt.Printf("Total time: %s\nAverage time between set and get: %s", totalTime, totalTime/N)

	cancel()
	time.Sleep(time.Nanosecond * 10000)

}
