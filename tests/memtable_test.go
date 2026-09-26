package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/KingrogKDR/omni/store"
	"github.com/KingrogKDR/omni/store/memtable"
)

func testStores() []struct {
	name string
	new  func(*testing.T) (store.Store, error)
} {
	return []struct {
		name string
		new  func(*testing.T) (store.Store, error)
	}{
		{
			name: "sorted slice",
			new: func(t *testing.T) (store.Store, error) {
				t.Helper()
				return memtable.NewMemtable(memtable.SORTED_SLICE)
			},
		},
		{
			name: "balanced tree",
			new: func(t *testing.T) (store.Store, error) {
				t.Helper()
				return memtable.NewMemtable(memtable.BALANCED_TREE)
			},
		},
		{
			name: "skip list",
			new: func(t *testing.T) (store.Store, error) {
				t.Helper()
				return memtable.NewMemtable(memtable.SKIPLIST)
			},
		},
	}
}

func TestPutGet(t *testing.T) {
	for _, tt := range testStores() {
		t.Run(tt.name, func(t *testing.T) {
			s, err := tt.new(t)
			if err != nil {
				t.Fatal(err)
			}

			key := []byte("hello")
			value := []byte("world")

			if err := s.Put(store.DefaultCF, key, value); err != nil {
				t.Fatalf("Put() error = %v", err)
			}

			got, found, err := s.Get(store.DefaultCF, key)
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}

			if !found {
				t.Fatal("Get() found = false, want true")
			}

			if string(got) != string(value) {
				t.Fatalf("Get() = %q, want %q", got, value)
			}
		})
	}
}

func TestGetMissing(t *testing.T) {
	for _, tt := range testStores() {
		t.Run(tt.name, func(t *testing.T) {
			s, err := tt.new(t)
			if err != nil {
				t.Fatal(err)
			}

			got, found, err := s.Get(
				store.DefaultCF,
				[]byte("does-not-exist"),
			)

			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}

			if found {
				t.Fatal("Get() found = true, want false")
			}

			if got != nil {
				t.Fatalf("Get() value = %q, want nil", got)
			}
		})
	}
}

func TestOverwrite(t *testing.T) {
	for _, tt := range testStores() {
		t.Run(tt.name, func(t *testing.T) {
			s, err := tt.new(t)
			if err != nil {
				t.Fatal(err)
			}

			key := []byte("foo")

			if err := s.Put(store.DefaultCF, key, []byte("one")); err != nil {
				t.Fatal(err)
			}

			if err := s.Put(store.DefaultCF, key, []byte("two")); err != nil {
				t.Fatal(err)
			}

			got, found, err := s.Get(store.DefaultCF, key)
			if err != nil {
				t.Fatal(err)
			}

			if !found {
				t.Fatal("key should exist")
			}

			if string(got) != "two" {
				t.Fatalf("got %q, want %q", got, "two")
			}
		})
	}
}

func TestDeleteThenGet(t *testing.T) {
	for _, tt := range testStores() {
		t.Run(tt.name, func(t *testing.T) {
			s, err := tt.new(t)
			if err != nil {
				t.Fatal(err)
			}

			key := []byte("foo")

			if err := s.Put(store.DefaultCF, key, []byte("bar")); err != nil {
				t.Fatal(err)
			}

			if err := s.Delete(store.DefaultCF, key); err != nil {
				t.Fatal(err)
			}

			got, found, err := s.Get(store.DefaultCF, key)
			if err != nil {
				t.Fatal(err)
			}

			if !errors.Is(err, store.ErrKeyNotFound) {
				t.Fatalf("Get() error = %v, want %v", err, store.ErrKeyNotFound)
			}

			if found {
				t.Fatal("deleted key was found")
			}

			if got != nil {
				t.Fatalf("Get() value = %q, want nil", got)
			}
		})
	}
}

func TestTTLExpiry(t *testing.T) {
	for _, tt := range testStores() {
		t.Run(tt.name, func(t *testing.T) {
			s, err := tt.new(t)
			if err != nil {
				t.Fatal(err)
			}

			key := []byte("temporary")

			if err := s.PutWithTTL(
				store.DefaultCF,
				key,
				[]byte("label"),
				50*time.Millisecond,
			); err != nil {
				t.Fatal(err)
			}

			value, found, err := s.Get(store.DefaultCF, key)
			if err != nil {
				t.Fatalf("Get() before expiration: %v", err)
			}

			if !found {
				t.Fatal("value should exist before expiration")
			}

			if string(value) != "label" {
				t.Fatalf("value = %q, want %q", value, "label")
			}

			time.Sleep(100 * time.Millisecond)

			value, found, err = s.Get(store.DefaultCF, key)

			if !errors.Is(err, store.ErrKeyNotFound) {
				t.Fatalf("Get() error = %v, want %v", err, store.ErrKeyNotFound)
			}

			if found {
				t.Fatal("expired value was found")
			}

			if value != nil {
				t.Fatalf("Get() value = %q, want nil", value)
			}
		})
	}
}
