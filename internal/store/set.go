package store

import "sync"

type SetCommand struct {
	Key, Val string
}

func (c *SetCommand) Execute(storeCtx *sync.Map) (string, error) {
	storeCtx.Store(c.Key, c.Val)
	return c.Val, nil
}
