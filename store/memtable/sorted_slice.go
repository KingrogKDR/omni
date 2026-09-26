package memtable

import (
	"bytes"
	"sort"
	"sync"

	"github.com/KingrogKDR/omni/store"
)

type item struct {
	key   []byte
	entry store.Entry
}

type SortedSlice struct {
	mu sync.RWMutex
	ds []item
}

func NewSortedSlice() *SortedSlice {
	return &SortedSlice{
		ds: make([]item, 0),
	}
}

func (l *SortedSlice) Insert(key []byte, e *store.Entry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	i := sort.Search(len(l.ds), func(i int) bool {
		return bytes.Compare(l.ds[i].key, key) >= 0
	})

	l.ds = append(l.ds, item{})

	copy(l.ds[i+1:], l.ds[i:])

	l.ds[i].key = bytes.Clone(key)
	l.ds[i].entry = *e
	return nil
}

func (l *SortedSlice) Delete(key []byte) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	// here, search looks for the first index where a condition becomes true instead of
	// looking for an exact match. So, bytes.Compare is used.
	i := sort.Search(len(l.ds), func(i int) bool {
		return bytes.Compare(l.ds[i].key, key) >= 0
	})

	if i == len(l.ds) || !bytes.Equal(l.ds[i].key, key) {
		return store.ErrKeyNotFound
	}

	l.ds[i].entry.Tombstone = true

	return nil
}

func (l *SortedSlice) Get(key []byte) (store.Entry, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Same reason. See the comment in Delete
	i := sort.Search(len(l.ds), func(i int) bool {
		return bytes.Compare(l.ds[i].key, key) >= 0
	})

	if i == len(l.ds) || !bytes.Equal(l.ds[i].key, key) {
		return store.Entry{}, store.ErrKeyNotFound
	}

	return l.ds[i].entry, nil
}
