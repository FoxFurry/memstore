package command

import "github.com/google/btree"

type GetCommand struct {
	Key string
}

func (c *GetCommand) Execute(storage *btree.BTree) (string, error) {
	return storage.Get(btree.).(string), nil
}
