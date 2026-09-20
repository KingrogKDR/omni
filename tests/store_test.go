package tests

// import (
// 	"errors"
// 	"fmt"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/KingrogKDR/omni/store"
// )

// // newTestStore is the single point of change when a later iteration wants
// // to re-run this exact suite against a new engine (WAL-backed, SSTable-
// // backed, full LSM). Swap this constructor; every test below stays as-is.
// func newTestStore(t *testing.T) store.Store {
// 	t.Helper()
// 	return NewMemtable()
// }

// func TestPutGet(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	if err := s.Put(store.DefaultCF, []byte("hello"), []byte("world")); err != nil {
// 		t.Fatalf("Put failed: %v", err)
// 	}

// 	val, found, err := s.Get(store.DefaultCF, []byte("hello"))
// 	if err != nil {
// 		t.Fatalf("Get failed: %v", err)
// 	}
// 	if !found {
// 		t.Fatalf("expected key to be found")
// 	}
// 	if string(val) != "world" {
// 		t.Fatalf("got %q, want %q", val, "world")
// 	}
// }

// func TestGetMissing(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	_, found, err := s.Get(store.DefaultCF, []byte("nope"))
// 	if err != nil {
// 		t.Fatalf("Get failed: %v", err)
// 	}
// 	if found {
// 		t.Fatalf("expected key to be absent")
// 	}
// }

// func TestOverwrite(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	s.Put(store.DefaultCF, []byte("k"), []byte("v1"))
// 	s.Put(store.DefaultCF, []byte("k"), []byte("v2"))

// 	val, found, _ := s.Get(store.DefaultCF, []byte("k"))
// 	if !found {
// 		t.Fatalf("expected key to be found")
// 	}
// 	if string(val) != "v2" {
// 		t.Fatalf("got %q, want %q (overwrite should win)", val, "v2")
// 	}
// }

// func TestDeleteThenGet(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	s.Put(store.DefaultCF, []byte("k"), []byte("v"))
// 	if err := s.Delete(store.DefaultCF, []byte("k")); err != nil {
// 		t.Fatalf("Delete failed: %v", err)
// 	}

// 	_, found, _ := s.Get(store.DefaultCF, []byte("k"))
// 	if found {
// 		t.Fatalf("expected key to be gone after delete (tombstone)")
// 	}
// }

// func TestDeleteMissingKeyIsNotError(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	if err := s.Delete(store.DefaultCF, []byte("never-existed")); err != nil {
// 		t.Fatalf("Delete on missing key should not error, got: %v", err)
// 	}
// }

// func TestReviveAfterDelete(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	s.Put(store.DefaultCF, []byte("k"), []byte("v1"))
// 	s.Delete(store.DefaultCF, []byte("k"))
// 	s.Put(store.DefaultCF, []byte("k"), []byte("v2")) // re-write after tombstone

// 	val, found, _ := s.Get(store.DefaultCF, []byte("k"))
// 	if !found {
// 		t.Fatalf("expected key to be found after re-write")
// 	}
// 	if string(val) != "v2" {
// 		t.Fatalf("got %q, want %q", val, "v2")
// 	}
// }

// func TestTTLExpiry(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	if err := s.PutWithTTL(store.DefaultCF, []byte("k"), []byte("v"), 20*time.Millisecond); err != nil {
// 		t.Fatalf("PutWithTTL failed: %v", err)
// 	}

// 	_, found, _ := s.Get(store.DefaultCF, []byte("k"))
// 	if !found {
// 		t.Fatalf("expected key to be found before expiry")
// 	}

// 	time.Sleep(40 * time.Millisecond)

// 	_, found, _ = s.Get(store.DefaultCF, []byte("k"))
// 	if found {
// 		t.Fatalf("expected key to be expired")
// 	}
// }

// func TestTTLDoesNotAffectOtherKeys(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	s.PutWithTTL(store.DefaultCF, []byte("short"), []byte("v"), 10*time.Millisecond)
// 	s.Put(store.DefaultCF, []byte("forever"), []byte("v"))

// 	time.Sleep(30 * time.Millisecond)

// 	_, found, _ := s.Get(store.DefaultCF, []byte("short"))
// 	if found {
// 		t.Fatalf("expected short-TTL key to be expired")
// 	}
// 	_, found, _ = s.Get(store.DefaultCF, []byte("forever"))
// 	if !found {
// 		t.Fatalf("expected no-TTL key to still be present")
// 	}
// }

// func TestManyKeysOrderingIndependence(t *testing.T) {
// 	// The Store contract makes no ordering guarantees (Iteration 0: no
// 	// scans exposed). This test only checks that insertion order doesn't
// 	// affect correctness of point lookups, not that any ordering exists.
// 	s := newTestStore(t)
// 	defer s.Close()

// 	n := 500
// 	for i := n - 1; i >= 0; i-- { // insert in reverse
// 		key := []byte(fmt.Sprintf("key-%04d", i))
// 		val := []byte(fmt.Sprintf("val-%04d", i))
// 		if err := s.Put(store.DefaultCF, key, val); err != nil {
// 			t.Fatalf("Put failed at i=%d: %v", i, err)
// 		}
// 	}

// 	for i := 0; i < n; i++ {
// 		key := []byte(fmt.Sprintf("key-%04d", i))
// 		want := fmt.Sprintf("val-%04d", i)
// 		val, found, _ := s.Get(store.DefaultCF, key)
// 		if !found {
// 			t.Fatalf("key-%04d missing", i)
// 		}
// 		if string(val) != want {
// 			t.Fatalf("key-%04d: got %q, want %q", i, val, want)
// 		}
// 	}
// }

// func TestConcurrentAccess(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	var wg sync.WaitGroup
// 	workers := 16
// 	opsPerWorker := 200

// 	for w := 0; w < workers; w++ {
// 		wg.Add(1)
// 		go func(w int) {
// 			defer wg.Done()
// 			for i := 0; i < opsPerWorker; i++ {
// 				key := []byte(fmt.Sprintf("w%d-k%d", w, i))
// 				val := []byte(fmt.Sprintf("w%d-v%d", w, i))
// 				if err := s.Put(store.DefaultCF, key, val); err != nil {
// 					t.Errorf("Put failed: %v", err)
// 					return
// 				}
// 				got, found, err := s.Get(store.DefaultCF, key)
// 				if err != nil {
// 					t.Errorf("Get failed: %v", err)
// 					return
// 				}
// 				if !found || string(got) != string(val) {
// 					t.Errorf("worker %d: got %q found=%v, want %q", w, got, found, val)
// 					return
// 				}
// 			}
// 		}(w)
// 	}
// 	wg.Wait()
// }

// func TestEmptyValue(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	if err := s.Put(store.DefaultCF, []byte("k"), []byte{}); err != nil {
// 		t.Fatalf("Put with empty value failed: %v", err)
// 	}
// 	val, found, _ := s.Get(store.DefaultCF, []byte("k"))
// 	if !found {
// 		t.Fatalf("expected key with empty value to be found")
// 	}
// 	if len(val) != 0 {
// 		t.Fatalf("expected empty value, got %q", val)
// 	}
// }

// // ---- Column family tests ----

// func TestDefaultCFExistsWithoutSetup(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	names := s.ColumnFamilies()
// 	if len(names) != 1 || names[0] != store.DefaultCF {
// 		t.Fatalf("expected only DefaultCF to exist initially, got %v", names)
// 	}
// }

// func TestCreateColumnFamily(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	if err := s.CreateColumnFamily("exact_cache"); err != nil {
// 		t.Fatalf("CreateColumnFamily failed: %v", err)
// 	}

// 	names := s.ColumnFamilies()
// 	if len(names) != 2 {
// 		t.Fatalf("expected 2 column families, got %v", names)
// 	}
// }

// func TestCreateColumnFamilyDuplicateErrors(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	if err := s.CreateColumnFamily("cf1"); err != nil {
// 		t.Fatalf("first create should succeed: %v", err)
// 	}
// 	err := s.CreateColumnFamily("cf1")
// 	if !errors.Is(err, store.ErrColumnFamilyExists) {
// 		t.Fatalf("expected ErrColumnFamilyExists, got %v", err)
// 	}
// }

// func TestOperationOnUnknownCFErrors(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	if _, _, err := s.Get("does-not-exist", []byte("k")); !errors.Is(err, store.ErrColumnFamilyNotFound) {
// 		t.Fatalf("Get: expected ErrColumnFamilyNotFound, got %v", err)
// 	}
// 	if err := s.Put("does-not-exist", []byte("k"), []byte("v")); !errors.Is(err, store.ErrColumnFamilyNotFound) {
// 		t.Fatalf("Put: expected ErrColumnFamilyNotFound, got %v", err)
// 	}
// 	if err := s.Delete("does-not-exist", []byte("k")); !errors.Is(err, store.ErrColumnFamilyNotFound) {
// 		t.Fatalf("Delete: expected ErrColumnFamilyNotFound, got %v", err)
// 	}
// }

// func TestColumnFamiliesAreIsolated(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	s.CreateColumnFamily("cf_a")
// 	s.CreateColumnFamily("cf_b")

// 	s.Put("cf_a", []byte("k"), []byte("value-in-a"))
// 	s.Put("cf_b", []byte("k"), []byte("value-in-b"))

// 	valA, found, _ := s.Get("cf_a", []byte("k"))
// 	if !found || string(valA) != "value-in-a" {
// 		t.Fatalf("cf_a: got %q found=%v, want %q", valA, found, "value-in-a")
// 	}

// 	valB, found, _ := s.Get("cf_b", []byte("k"))
// 	if !found || string(valB) != "value-in-b" {
// 		t.Fatalf("cf_b: got %q found=%v, want %q", valB, found, "value-in-b")
// 	}

// 	// Same key must not exist in DefaultCF just because it exists elsewhere.
// 	_, found, _ = s.Get(store.DefaultCF, []byte("k"))
// 	if found {
// 		t.Fatalf("expected key absent from DefaultCF")
// 	}
// }

// // ---- WriteBatch tests ----

// func TestWriteBatchAppliesAllOps(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	s.CreateColumnFamily("cf_a")

// 	err := s.WriteBatch([]Op{
// 		{CF: store.DefaultCF, Key: []byte("k1"), Value: []byte("v1")},
// 		{CF: "cf_a", Key: []byte("k2"), Value: []byte("v2")},
// 		{CF: store.DefaultCF, Key: []byte("k3"), Value: []byte("v3"), HasTTL: true, TTL: time.Hour},
// 	})
// 	if err != nil {
// 		t.Fatalf("WriteBatch failed: %v", err)
// 	}

// 	v1, found, _ := s.Get(store.DefaultCF, []byte("k1"))
// 	if !found || string(v1) != "v1" {
// 		t.Fatalf("k1: got %q found=%v", v1, found)
// 	}
// 	v2, found, _ := s.Get("cf_a", []byte("k2"))
// 	if !found || string(v2) != "v2" {
// 		t.Fatalf("k2: got %q found=%v", v2, found)
// 	}
// 	v3, found, _ := s.Get(store.DefaultCF, []byte("k3"))
// 	if !found || string(v3) != "v3" {
// 		t.Fatalf("k3: got %q found=%v", v3, found)
// 	}
// }

// func TestWriteBatchWithDelete(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	s.Put(store.DefaultCF, []byte("k"), []byte("v"))

// 	err := s.WriteBatch([]Op{
// 		{CF: store.DefaultCF, Key: []byte("k"), Delete: true},
// 		{CF: store.DefaultCF, Key: []byte("other"), Value: []byte("v2")},
// 	})
// 	if err != nil {
// 		t.Fatalf("WriteBatch failed: %v", err)
// 	}

// 	_, found, _ := s.Get(store.DefaultCF, []byte("k"))
// 	if found {
// 		t.Fatalf("expected k to be deleted by batch")
// 	}
// 	v, found, _ := s.Get(store.DefaultCF, []byte("other"))
// 	if !found || string(v) != "v2" {
// 		t.Fatalf("other: got %q found=%v", v, found)
// 	}
// }

// // TestWriteBatchAllOrNothingOnUnknownCF is the core atomicity guarantee:
// // if any Op in the batch references a CF that doesn't exist, none of the
// // batch's writes should be visible afterward — not just the failing one.
// func TestWriteBatchAllOrNothingOnUnknownCF(t *testing.T) {
// 	s := newTestStore(t)
// 	defer s.Close()

// 	err := s.WriteBatch([]Op{
// 		{CF: store.DefaultCF, Key: []byte("should-not-land"), Value: []byte("v1")},
// 		{CF: "no-such-cf", Key: []byte("k2"), Value: []byte("v2")},
// 	})
// 	if !errors.Is(err, store.ErrColumnFamilyNotFound) {
// 		t.Fatalf("expected ErrColumnFamilyNotFound, got %v", err)
// 	}

// 	_, found, _ := s.Get(store.DefaultCF, []byte("should-not-land"))
// 	if found {
// 		t.Fatalf("batch partially applied: 'should-not-land' should not be visible " +
// 			"since the batch as a whole failed validation")
// 	}
// }

// func TestWriteBatchConcurrentWithReads(t *testing.T) {
// 	// A reader should never observe a batch half-applied: both keys appear
// 	// together, or neither does.
// 	s := newTestStore(t)
// 	defer s.Close()

// 	var wg sync.WaitGroup
// 	iterations := 200

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		for i := 0; i < iterations; i++ {
// 			s.WriteBatch([]Op{
// 				{CF: store.DefaultCF, Key: []byte("pair-a"), Value: []byte(fmt.Sprintf("v%d", i))},
// 				{CF: store.DefaultCF, Key: []byte("pair-b"), Value: []byte(fmt.Sprintf("v%d", i))},
// 			})
// 		}
// 	}()

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		for i := 0; i < iterations; i++ {
// 			va, foundA, _ := s.Get(store.DefaultCF, []byte("pair-a"))
// 			vb, foundB, _ := s.Get(store.DefaultCF, []byte("pair-b"))
// 			if foundA != foundB {
// 				t.Errorf("observed batch half-applied: pair-a found=%v, pair-b found=%v", foundA, foundB)
// 				return
// 			}
// 			if foundA && string(va) != string(vb) {
// 				t.Errorf("observed batch half-applied: pair-a=%q, pair-b=%q", va, vb)
// 				return
// 			}
// 		}
// 	}()

// 	wg.Wait()
// }
