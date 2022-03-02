package command

import (
	"github.com/google/btree"
)

type set struct {
	pair
}

func (cmd *set) Execute(storage *btree.BTree) (string, error) {
	storage.ReplaceOrInsert(cmd.pair)
	return cmd.value, nil
}

func (cmd *set) Type() CommandType {
	return Write
}

func Set(key, value string) Command {
	return &set{
		pair: pair{
			key:   key,
			value: value,
		},
	}
}
