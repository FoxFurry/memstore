package command

import (
	"github.com/google/btree"
	"strings"
)

type CommandType int

const (
	Read CommandType = iota
	Write
)

// Command interface describes how a command affects the low level storage. Each command should implement 3 methods:
// Execute(), Type() and Key()
//
// - Execute should represent how certain Command will affect the storage, Read commands should return value of specified key
// and Write commands should return it's own value as result. Any error returned will be considered as fatal (for now).
// Interaction with storage is done via btree.Btree
//
// - Type should return CommandType for a certain Command. If a Command is not writing any data - it is considered as read, otherwise - as write
//
// - Key should return single string which represents what Command is querying. Embed pair struct to automatically embed Key() method
type Command interface {
	Execute(storage *btree.BTree) (string, error)
	Type() CommandType
	Key() string
}

type pair struct {
	key   string
	value string
}

// Less implements btree interface
func (p pair) Less(b btree.Item) bool {
	return strings.Compare(p.key, b.(pair).key) < 0
}

// Key implements Command interface
func (p pair) Key() string {
	return p.key
}
