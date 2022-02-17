package main

import (
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

	cmds := []store.ICommand{
		&store.SetCommand{
			Key: "test",
			Val: "test",
		},
	}

	result, err := storage.Execute(cmds)

	fmt.Printf("Res: %v\n------\nErr: %v\n\n\n", result, err)
	cmds = []store.ICommand{
		&store.GetCommand{
			Key: "test",
		},
	}

	result, err = storage.Execute(cmds)
	fmt.Printf("Res: %v\n------\nErr: %v\n\n", result, err)

	cancel()
}
