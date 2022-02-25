package cluster

import (
	"KeyValueHTTPStore/internal/command"
	"context"
	"github.com/google/btree"
	"log"
	"strconv"
	"sync"
)

type inode interface {
	startJournal(ctx context.Context)

	execute(cmd command.Command) (string, error)
	journal() chan<- []command.Command
	snapshot() inode
	String() string
}

type node struct {
	cloneMutex     sync.RWMutex
	storageMutex   sync.RWMutex
	storage, clone *btree.BTree
	nodeID         int
	journalQueue   chan []command.Command
}

func newNode(ID int) inode {
	n := &node{
		nodeID:       ID,
		storage:      btree.New(32), // TODO: Why do I use 32?
		cloneMutex:   sync.RWMutex{},
		storageMutex: sync.RWMutex{},
		journalQueue: make(chan []command.Command, 50), // TODO: Why 50??? So many questions
	}
	n.clone = n.storage.Clone()

	return n
}

func (n *node) execute(cmd command.Command) (string, error) { // IMPORTANT: This function doesn't lock the mutex
	return cmd.Execute(n.storage)
}

func (n *node) journal() chan<- []command.Command {
	return n.journalQueue
}

func (n *node) snapshot() inode {
	n.cloneMutex.RLock()
	clone := n.clone.Clone()
	n.cloneMutex.RUnlock()

	return &node{
		storage:    clone,          // Lazy cow copy
		nodeID:     n.nodeID,       // Snapshot should have same ID
		cloneMutex: sync.RWMutex{}, // Each snapshot will have each own mutex
		clone:      nil,            // snapshot should not have any clones
	}
}

func (n *node) startJournal(ctx context.Context) {
	for {
		select {
		case block := <-n.journalQueue:

			for _, cmd := range block {
				_, err := cmd.Execute(n.storage)
				if err != nil {
					log.Panicf("Pizdec deadlock nafig: %v", err)
				}
			}

			n.cloneMutex.Lock()
			n.clone = n.storage.Clone()
			n.cloneMutex.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (n *node) String() string {
	return strconv.Itoa(n.nodeID)
}
