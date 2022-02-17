package store

import (
	"context"
	"sync"
)

type IStore interface {
	Execute(commandQueue []ICommand) ([]string, error)
	Start(ctx context.Context)
}

type KeyValueStore struct {
	store   Storage
	journal IJournal
}

func NewKeyValueStore() IStore {
	return &KeyValueStore{
		store: Storage{
			S:             sync.Map{},
			Snapshot:      sync.Map{},
			SnapshotMutex: sync.RWMutex{},
		},
		journal: NewJournal(),
	}
}

func (s *KeyValueStore) Start(ctx context.Context) {
	s.journal.Listen(ctx, &s.store)
}

func (s *KeyValueStore) Execute(commandQueue []ICommand) ([]string, error) {
	var allResults []string

	// We take snapshot of current storage and will test this commands on the snapshot
	storage := s.store.TakeSnapshot()

	for _, command := range commandQueue {
		// Since snapshot correctly represents storage current state - we can just return execution results on snapshot
		cmdResult, err := command.Execute(&storage)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, cmdResult)
	}

	// If no errors occurred - add commands to a real execute queue
	s.journal.AddCmdBlock(commandQueue)

	return allResults, nil
}
