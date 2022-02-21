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

type Command interface {
	Execute(storage *btree.BTree) (string, error)
	Type() CommandType
	Key() string
}

type pair struct {
	key   string
	value string
}

func (p pair) Less(b btree.Item) bool {
	return strings.Compare(p.key, b.(pair).key) < 0 // Fucking ugly
}
