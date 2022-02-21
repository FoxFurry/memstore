package main

import (
	"KeyValueHTTPStore/internal/cluster"
	"KeyValueHTTPStore/internal/command"
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	storage := cluster.New()
	storage.Initialize(ctx)

	time.Sleep(time.Microsecond * 200)

	cmds := []command.Command{
		command.Set("test", "testv"),
		command.Set("test312we", "testv"),
		command.Set("tesdsadast", "testv"),
		command.Set("test32", "testv"),
		command.Set("teasdast", "testv"),
		command.Set("@jmvj785n2", "testv2"),
	}

	result, err := storage.Execute(cmds)

	time.Sleep(time.Nanosecond * 150)

	fmt.Printf("Res: %v\n------\nErr: %v\n\n\n", result, err)

	cmds = []command.Command{
		command.Get("test"),
		command.Get("@jmvj785n2"),
	}

	result, err = storage.Execute(cmds)
	fmt.Printf("Res: %v\n------\nErr: %v\n\n", result, err)

	time.Sleep(time.Nanosecond * 10000)

	cancel()
}
