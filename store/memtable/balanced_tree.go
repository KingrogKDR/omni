package memtable

import (
	"sync"

	"github.com/KingrogKDR/omni/store"
)

type BalancedTree struct {
	mu sync.RWMutex
}

func NewBalancedTree() *BalancedTree {
	return &BalancedTree{}
}

func (b *BalancedTree) Insert(key []byte, e *store.Entry) error {
	return nil
}

func (b *BalancedTree) Delete(key []byte) error {
	return nil
}

func (b *BalancedTree) Get(key []byte) (store.Entry, error) {
	return store.Entry{}, nil
}
