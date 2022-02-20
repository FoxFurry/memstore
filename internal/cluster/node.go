package cluster

import (
	"context"
	"github.com/google/btree"
	"strconv"
	"sync"
)

type inode interface {
	execute(command interface{}) (string, error)
	addToJournal(commands []interface{})
	snapshot() inode
	startJournal(ctx context.Context)
	String() string
}

type node struct {
	storageMutex sync.RWMutex
	storage      *btree.BTree
	nodeID       int
	cmdQueue     chan []interface{}
}

func newNode(ctx context.Context, ID int) inode {
	n := &node{
		nodeID:       ID,
		storage:      btree.New(32), // TODO: Why do I use 32?
		storageMutex: sync.RWMutex{},
	}
	n.startJournal(ctx)
	return n
}

func (n *node) execute(command interface{}) (string, error) {
	return "", nil
}

func (n *node) snapshot() inode {
	n.storageMutex.RLock()
	defer n.storageMutex.RUnlock()

	return &node{
		storage:      n.storage.Clone(),
		nodeID:       n.nodeID,
		storageMutex: sync.RWMutex{}, // Each snapshot will have each own mutex
	}
}

func (n *node) addToJournal(cmds []interface{}) {
	n.cmdQueue <- cmds
}

func (n *node) startJournal(ctx context.Context) {
	for {
		select {
		case block := <-n.cmdQueue:
			n.storageMutex.Lock()

			for _, _ := range block {
				// execute

			}

			n.storageMutex.Unlock()
		}
	}
}

func (n *node) String() string {
	return strconv.Itoa(n.nodeID)
}
