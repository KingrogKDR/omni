package memtable

import (
	"sync"

	"github.com/KingrogKDR/omni/store"
)

type SortedSlice struct {
	mu sync.RWMutex
}

func NewSortedSlice() *SortedSlice {
	return &SortedSlice{}
}

func (l *SortedSlice) Insert(key []byte, e *store.Entry) error {
	return nil
}

func (l *SortedSlice) Delete(key []byte) error {
	return nil
}

func (l *SortedSlice) Get(key []byte) (store.Entry, error) {
	return store.Entry{}, nil
}
