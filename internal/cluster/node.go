package cluster

import (
	"context"
	"github.com/FoxFurry/memstore/internal/command"
	"github.com/google/btree"
	"log"
	"strconv"
	"sync"
)

// inode represents interface of a shard (node)
// TODO: research if using interfaces slows down execution
type inode interface {
	startQueueListener(ctx context.Context)

	execute(cmd command.Command) (string, error)
	addToQueue(block []command.Command)
	snapshot() inode
	String() string
}

type node struct {
	storageMutex sync.RWMutex
	storage      *btree.BTree
	nodeID       int
	nodeQueue    chan []command.Command
}

func newNode(ID int) inode {
	n := &node{
		nodeID:       ID,
		storage:      btree.New(4), // TODO: Why do I use 4?
		storageMutex: sync.RWMutex{},
		nodeQueue:    make(chan []command.Command, 50), // TODO: Why 50??? So many questions
	}

	return n
}

// execute passes storage to execute method of a command. It is thread-safe, but mutex here might slow down performance
// TODO: rethink mutex usage
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

// addToQueue adds transaction to queue
func (n *node) addToQueue(transaction []command.Command) {
	n.nodeQueue <- transaction
}

// snapshot creates a fast copy of a node itself. snapshot is thread-safe method
func (n *node) snapshot() inode {
	n.storageMutex.RLock()
	defer n.storageMutex.RUnlock()

	return &node{
		storage:      n.storage.Clone(), // Lazy cow copy
		nodeID:       n.nodeID,          // Snapshot should have same ID
		storageMutex: sync.RWMutex{},    // Each snapshot will have each own mutex
	}
}

// startQueueListener listens for incoming transactions and executes them
func (n *node) startQueueListener(ctx context.Context) {
	log.Printf("Starting queue listener %d\n", n.nodeID)
	for {
		select {
		case transaction := <-n.nodeQueue: // Get transaction from queue
			// TODO: Test per-transaction mutex vs per-command mutex
			for _, cmd := range transaction { // Go through every command in transaction
				_, err := cmd.Execute(n.storage) // Execute each command. We don't care about result here
				if err != nil {
					log.Panicf("%s: %s", errExecutionFailed, err)
				}
			}

		case <-ctx.Done():
			log.Printf("Stopping queue listener %d\n", n.nodeID)
			return
		}
	}
}

// String method is required by btree interface
func (n *node) String() string {
	return strconv.Itoa(n.nodeID)
}
