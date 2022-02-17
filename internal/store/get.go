package store

import (
	storeErr "KeyValueHTTPStore/internal/error"
	"sync"
)

type GetCommand struct {
	Key string
}

func (c *GetCommand) Execute(storeCtx *sync.Map) (string, error) {
	val, ok := storeCtx.Load(c.Key)
	if !ok {
		return "", storeErr.NewStoreKeyNotFound(c.Key)
	}

	strVal, ok := val.(string)
	if !ok {
		return "", storeErr.NewCouldNotCast(val)
	}

	return strVal, nil
}
