package command

import (
	"github.com/google/btree"
)

type set struct {
	pair
}

func (cmd *set) Execute(storage *btree.BTree) (string, error) {
	storage.ReplaceOrInsert(cmd.pair)
	return cmd.pair.value, nil
}

func (cmd *set) Key() string {
	return cmd.key
}

func Set(key, value string) Command {
	return &set{
		pair: pair{
			key:   key,
			value: value,
		},
	}
}
