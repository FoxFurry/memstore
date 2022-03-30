package main

import (
	"context"
	"github.com/FoxFurry/GoKeyValueStore/internal/http/server"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	ctx, _ := context.WithCancel(context.Background())

	srv := server.New(ctx)

	srv.Start()
}
