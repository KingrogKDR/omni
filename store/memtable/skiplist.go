package memtable

import (
	"sync"

	"github.com/KingrogKDR/omni/store"
)

type SkipList struct {
	mu sync.RWMutex
}

func NewSkipList() *SkipList {
	return &SkipList{}
}

func (l *SkipList) Insert(key []byte, e *store.Entry) error {
	return nil
}

func (l *SkipList) Delete(key []byte) error {
	return nil
}

func (l *SkipList) Get(key []byte) (store.Entry, error) {
	return store.Entry{}, nil
}
