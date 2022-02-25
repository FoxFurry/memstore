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
	ctx, cancel := context.WithCancel(context.Background())
	storage := cluster.New()
	storage.Initialize(ctx)

	cmds := []command.Command{
		command.Set("test", "testv"),
	}

	result, err := storage.Execute(cmds)

	time.Sleep(time.Nanosecond * 500)

	fmt.Printf("Res: %v\n------\nErr: %v\n\n\n", result, err)

	cmds = []command.Command{
		command.Get("test"),
	}

	result, err = storage.Execute(cmds)
	fmt.Printf("Res: %v\n------\nErr: %v\n\n", result, err)

	time.Sleep(time.Nanosecond * 10000)

	cancel()
}
