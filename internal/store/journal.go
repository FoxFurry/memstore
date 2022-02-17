package store

import (
	"context"
	"log"
	"sync"
)

type IJournal interface {
	AddCmdBlock([]ICommand)
	Listen(ctx context.Context, store *Storage)
}

type Journal struct {
	blocks chan []ICommand
}

func NewJournal() IJournal {
	return &Journal{
		blocks: make(chan []ICommand, 50),
	}
}

func (j *Journal) AddCmdBlock(block []ICommand) {
	j.blocks <- block
}

func (j *Journal) Listen(ctx context.Context, store *Storage) {
	log.Printf("Starting journal listener")

	for {
		select {
		case block := <-j.blocks:
			store.SnapshotMutex.Lock()

			for _, cmd := range block {
				cmd.Execute(&store.S)
			}

			store.Snapshot = copySyncMap(store.S)
			store.SnapshotMutex.Unlock()
		case <-ctx.Done():
			log.Printf("Stopping journal listener")
			return
		}
	}
}

func copySyncMap(m sync.Map) sync.Map {
	var cp sync.Map

	m.Range(func(k, v interface{}) bool {
		vm, ok := v.(sync.Map)
		if ok {
			cp.Store(k, copySyncMap(vm))
		} else {
			cp.Store(k, v)
		}

		return true
	})

	return cp
}
