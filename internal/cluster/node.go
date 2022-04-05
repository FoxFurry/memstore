package cluster

import (
	"context"
	"github.com/FoxFurry/memstore/internal/command"
	"github.com/google/btree"
	"log"
	"strconv"
	"sync"
)

type inode interface {
	startJournal(ctx context.Context)

	execute(cmd command.Command) (string, error)
	addToJournal(block []command.Command)
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
		storage:      btree.New(4), // TODO: Why do I use 32?
		storageMutex: sync.RWMutex{},
		journalQueue: make(chan []command.Command, 50), // TODO: Why 50??? So many questions
	}

	return n
}

func (n *node) execute(cmd command.Command) (string, error) {
	if cmd.Type() == command.Write {
		n.storageMutex.Lock()
		defer n.storageMutex.Unlock()
	} else {
		n.storageMutex.RLock()
		defer n.storageMutex.RUnlock()
	}

	return cmd.Execute(n.storage)
}

func (n *node) addToJournal(block []command.Command) {
	n.journalQueue <- block
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
	for {
		select {
		case block := <-n.journalQueue:

			for _, cmd := range block {
				_, err := cmd.Execute(n.storage)
				if err != nil {
					log.Panicf("%s: %s", errExecutionFailed, err)
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func (n *node) String() string {
	return strconv.Itoa(n.nodeID)
}
