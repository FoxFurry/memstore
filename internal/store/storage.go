package store

import "sync"

type Storage struct {
	S             sync.Map
	Snapshot      sync.Map
	SnapshotMutex sync.RWMutex
}

// TakeSnapshot returns a copy
func (s *Storage) TakeSnapshot() sync.Map {
	s.SnapshotMutex.RLock()
	defer s.SnapshotMutex.RUnlock()
	return s.Snapshot
}
