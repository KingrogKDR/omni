package kvstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type RecordType uint8

const (
	RecordPut RecordType = iota
	RecordDelete
)

var ErrKeyNotFound = errors.New("key not found")

const defaultDataDir = ".omni"

// record represents a single operation (for internal use)
type record struct {
	Type RecordType `json:"type"`
	Key  string     `json:"key"`
	Val  []byte     `json:"value,omitempty"` // only for PUT
}

// Config holds optional database settings (can be used to override the defaults)
type Config struct {
	DataDir   string
	InMemOnly bool // 0 means default persistence, 1 means in-memory map
}

func DefaultConfig() Config {
	return Config{
		DataDir:   "",
		InMemOnly: false,
	}
}

type Backend interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	View() error
	Close() error
}

type MemoryStore struct {
	inMemMap map[string][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		inMemMap: make(map[string][]byte),
	}
}

// store a copy of the value in memory
// avoids cases like this:
// 	buf := []byte("hello")
// 	db.Put("a", buf)
// 	buf[0] = 'X' -> becomes Xello

func (ms *MemoryStore) Put(key string, val []byte) error {
	copy := append([]byte(nil), val...)
	ms.inMemMap[key] = copy
	return nil
}

// similarly gets a copy of the value in memory
func (ms *MemoryStore) Get(key string) ([]byte, error) {
	value, ok := ms.inMemMap[key]
	if !ok {
		return nil, ErrKeyNotFound
	}
	copy := append([]byte(nil), value...)
	return copy, nil
}

func (ms *MemoryStore) Delete(key string) error {
	if _, ok := ms.inMemMap[key]; !ok {
		return ErrKeyNotFound
	}

	delete(ms.inMemMap, key)
	return nil
}

func (ms *MemoryStore) View() error {
	for key, value := range ms.inMemMap {
		fmt.Printf("%s: %s\n", key, value)
	}
	return nil
}

func (ms *MemoryStore) Close() error {
	return nil
}

type UpperStore struct {
	UpperDir string
}

func NewUpperStore(dir string) *UpperStore {
	return &UpperStore{
		UpperDir: dir,
	}
}

type LowerStore struct {
	LowerDir string
}

func NewLowerStore(dir string) *LowerStore {
	return &LowerStore{
		LowerDir: dir,
	}
}

type PersistentStore struct {
	dataDir string
	upper   *UpperStore
	lower   *LowerStore
}

func NewPersistentStore(dataDir string) (*PersistentStore, error) {
	upperDir := filepath.Join(dataDir, ".upper")
	lowerDir := filepath.Join(dataDir, ".lower")

	if err := os.MkdirAll(upperDir, 0o755); err != nil {
		return nil, fmt.Errorf("upper dir: %w", err)
	}

	if err := os.MkdirAll(lowerDir, 0o755); err != nil {
		return nil, fmt.Errorf("lower dir: %w", err)
	}

	return &PersistentStore{
		dataDir: dataDir,
		upper:   NewUpperStore(upperDir),
		lower:   NewLowerStore(lowerDir),
	}, nil
}

func (ps *PersistentStore) Put(key string, val []byte) error {
	rec := record{
		Type: RecordPut,
		Key:  key,
		Val:  val,
	}

	filePath := fmt.Sprintf("%s/segment1.log", ps.upper.UpperDir)

	fd, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer fd.Close()

	recBytes, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	recBytes = append(recBytes, '\n')
	_, err = fd.Write(recBytes)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	fmt.Println("OK")
	return nil
}

func (ps *PersistentStore) Get(key string) ([]byte, error) {
	filePath := fmt.Sprintf("%s/segment1.log", ps.upper.UpperDir)
	fd, err := os.OpenFile(filePath, os.O_RDONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer fd.Close()

	// it should not be readall, instead it should read from the bottom of the file and give the latest instance from there
	data, err := io.ReadAll(fd)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	if len(data) == 0 {
		return nil, ErrKeyNotFound
	}

	// var entries []record
	// if err := json.Unmarshal(data, &entries); err != nil {
	// 	return nil, fmt.Errorf("unmarshal JSON: %w", err)
	// }

	// value, exists := entries[key]
	// if !exists {
	// 	return nil, ErrKeyNotFound
	// }

	return nil, nil
}

func (ps *PersistentStore) Delete(key string) error {
	return nil
}

func (ps *PersistentStore) Close() error {
	return nil
}

func (ps *PersistentStore) View() error {
	// prints upper and lower in a beautiful format
	return nil
}
