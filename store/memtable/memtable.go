package memtable

import (
	"fmt"
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/KingrogKDR/omni/store"
)

type dataStructureType uint8

const (
	SORTED_SLICE dataStructureType = iota
	BALANCED_TREE
	SKIPLIST
)

type DataStructure interface {
	Insert(key []byte, e *store.Entry) error
	Delete(key []byte) error
	Get(key []byte) (store.Entry, error)
}

func newDataStructure(t dataStructureType) (DataStructure, error) {
	switch t {
	case SORTED_SLICE:
		return NewSortedSlice(), nil
	case BALANCED_TREE:
		return NewBalancedTree(), nil
	case SKIPLIST:
		return NewSkipList(), nil
	default:
		return nil, fmt.Errorf("unsupported data structure type: %v", t)
	}
}

type Memtable struct {
	mu sync.RWMutex

	// type of data structure used by the memTable
	typ dataStructureType

	// a map of all column families and their respective data structures
	cfs map[store.CFName]DataStructure
}

func NewMemtable(typ dataStructureType) (*Memtable, error) {
	ds, err := newDataStructure(typ)
	if err != nil {
		return nil, fmt.Errorf("creating data structure error: %w", err)
	}
	return &Memtable{
		cfs: map[store.CFName]DataStructure{store.DefaultCF: ds},
		typ: typ,
	}, nil
}

// -- column family operations

func (m *Memtable) CreateColumnFamily(cf store.CFName) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.cfs[cf]; exists {
		return fmt.Errorf("%w: %s", store.ErrColumnFamilyExists, cf)
	}
	dataStruct, err := newDataStructure(m.typ)
	if err != nil {
		return fmt.Errorf("create data structure: %w", err)
	}

	m.cfs[cf] = dataStruct
	return nil
}

func (m *Memtable) ColumnFamilies() []store.CFName {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return slices.Collect(maps.Keys(m.cfs))
}

func (m *Memtable) getDataStructure(
	cf store.CFName,
) (DataStructure, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ds, exists := m.cfs[cf]
	if !exists {
		return nil, fmt.Errorf(
			"%w: %s",
			store.ErrColumnFamilyNotFound,
			cf,
		)
	}

	return ds, nil
}

// -- read/write operations

func (m *Memtable) Put(cf store.CFName, key, value []byte) error {
	return m.putInternal(cf, key, value, store.Entry{})
}

func (m *Memtable) putInternal(cf store.CFName, key, value []byte, e store.Entry) error {
	ds, err := m.getDataStructure(cf)
	if err != nil {
		return err
	}
	e.Val = value
	if err := ds.Insert(key, &e); err != nil {
		return fmt.Errorf("insertion in data structure: %w", err)
	}

	return nil
}

func (m *Memtable) PutWithTTL(cf store.CFName, key, value []byte, ttl time.Duration) error {
	return m.putInternal(cf, key, value, store.Entry{
		HasExpiry: true,
		ExpiresAt: time.Now().Add(ttl),
	})
}

func (m *Memtable) Get(cf store.CFName, key []byte) (value []byte, found bool, err error) {
	now := time.Now()
	ds, err := m.getDataStructure(cf)
	if err != nil {
		return nil, false, err
	}

	entry, err := ds.Get(key)
	if err != nil {
		return nil, false, fmt.Errorf("get from data structure: %w", err)
	}

	if entry.Tombstone || entry.IsExpired(now) {
		return nil, false, store.ErrKeyNotFound
	}

	return entry.Val, true, nil
}

func (m *Memtable) Delete(cf store.CFName, key []byte) error {
	ds, err := m.getDataStructure(cf)
	if err != nil {
		return err
	}

	if err := ds.Delete(key); err != nil {
		return fmt.Errorf("delete from data structure: %w", err)
	}

	return nil
}

func (m *Memtable) Close() error {
	return nil
}
