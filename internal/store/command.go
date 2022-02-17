package store

import "sync"

type ICommand interface {
	Execute(storeCtx *sync.Map) (string, error)
}
