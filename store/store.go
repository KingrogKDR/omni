package store

import (
	"errors"
	"time"
)

// ErrKeyNotFound is returned by Get when a key is absent or has expired.
// Note: Get's second return value (bool) is the primary "found" signal;
// this error is reserved for implementations where the lookup
// itself can fail (e.g. disk I/O errors).
var ErrKeyNotFound = errors.New("store: key not found")

// ErrColumnFamilyNotFound is returned by any operation referencing a CFName
// that has not been created via CreateColumnFamily.
var ErrColumnFamilyNotFound = errors.New("store: column family not found")

// ErrColumnFamilyExists is returned by CreateColumnFamily when the name is
// already in use.
var ErrColumnFamilyExists = errors.New("store: column family already exists")

type CFName string

// DefaultCF is the default CF name created automatically by every Store implementation
const DefaultCF CFName = "default"

// Store is the public contract every engine implementation satisfies.
type Store interface {
	// CreateColumnFamily registers a new, empty column family. Returns
	// ErrColumnFamilyExists if the name is already registered.
	CreateColumnFamily(cf CFName) error

	// ColumnFamilies returns the names of all registered column families,
	// including DefaultCF.
	ColumnFamilies() []CFName

	// Put writes key -> value with no expiry, in the given column family. Overwrites existing key with new value.
	// Returns ErrColumnFamilyNotFound if cf has not been created.
	Put(cf CFName, key, value []byte) error

	// PutWithTTL writes key -> value that expires after ttl elapses. Overwrites existing key with new value.
	PutWithTTL(cf CFName, key, value []byte, ttl time.Duration) error

	// Get returns the value for key in the given column family. found is
	// false if the key was never written, was deleted, or has expired.
	Get(cf CFName, key []byte) (value []byte, found bool, err error)

	// Delete writes a tombstone for key. A subsequent Get returns found=false. Delete does not error if the key did not previously exist.
	Delete(cf CFName, key []byte) error

	// Close releases any resources held by the store. In this iteration
	// (pure in-memory) it is a no-op, but every implementation must provide
	// it since later iterations (WAL file handles, SSTable files) need it.
	Close() error
}

type EntryType uint8

const (
	PUT EntryType = iota + 1
	DELETE
)

type Entry struct {
	HasExpiry bool
	ExpiresAt time.Time
	Tombstone bool
	Val       []byte
	Typ       EntryType
}

func (e *Entry) IsExpired(now time.Time) bool {
	return e.HasExpiry && now.After(e.ExpiresAt)
}
