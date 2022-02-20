package command

import (
	"github.com/google/btree"
	"strings"
)

type ICommand interface {
	Execute(storage *btree.BTree) (string, error)
}

type Key string

func (k Key) Less(b btree.Item) bool {
	return strings.Compare(string(k), string(b.(Key))) > 0 // Fucking ugly
}
