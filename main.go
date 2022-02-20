package main

import (
	"KeyValueHTTPStore/internal/command"
	"KeyValueHTTPStore/internal/store"
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	storage := store.NewKeyValueStore()
	go storage.Start(ctx)
	time.Sleep(time.Second / 2)

	cmds := []command.ICommand{
		&command.SetCommand{
			Key: "test",
			Val: "test",
		},
	}

	result, err := storage.Execute(cmds)

	fmt.Printf("Res: %v\n------\nErr: %v\n\n\n", result, err)

	time.Sleep(time.Second)
	cmds = []command.ICommand{
		&command.GetCommand{
			Key: "test",
		},
	}

	result, err = storage.Execute(cmds)
	fmt.Printf("Res: %v\n------\nErr: %v\n\n", result, err)

	cancel()
}
