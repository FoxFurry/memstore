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
	storageMutex sync.RWMutex
	storage      *btree.BTree
	nodeID       int
	journalQueue chan []command.Command
}

func newNode(ID int) inode {
	n := &node{
		nodeID:       ID,
		storage:      btree.New(32), // TODO: Why do I use 32?
		storageMutex: sync.RWMutex{},
		journalQueue: make(chan []command.Command, 50), // TODO: Why 50??? So many questions
	}
	return n
}

func (n *node) execute(cmd command.Command) (string, error) { // IMPORTANT: This function doesn't lock the mutex
	return cmd.Execute(n.storage)
}

func (n *node) journal() chan<- []command.Command {
	return n.journalQueue
}

func (n *node) snapshot() inode {
	n.storageMutex.RLock()
	defer n.storageMutex.RUnlock()

	return &node{
		storage:      n.storage.Clone(), // Lazy cow copy
		nodeID:       n.nodeID,          // Snapshot should have same ID
		storageMutex: sync.RWMutex{},    // Each snapshot will have each own mutex
	}
}

func (n *node) startJournal(ctx context.Context) {
	log.Printf("Starting journal #%d", n.nodeID)

	for {
		select {
		case block := <-n.journalQueue:
			log.Printf("Journal #%d executing command block", n.nodeID)

			n.storageMutex.Lock() // Cannot take snapshots during block execution

			for _, cmd := range block {
				_, err := n.execute(cmd)
				if err != nil {
					log.Panicf("Pizdec deadlock nafig: %v", err)
				}
			}

			n.storageMutex.Unlock()
		case <-ctx.Done():
			log.Printf("Closing journal #%d", n.nodeID)
		}
	}
}

func (n *node) String() string {
	return strconv.Itoa(n.nodeID)
}
