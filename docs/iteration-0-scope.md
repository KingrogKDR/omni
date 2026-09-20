# Iteration 0 — Scope & Contracts

## API surface (locked for the whole project)

```go
Put(cf CFName, key, value []byte) error
PutWithTTL(cf CFName, key, value []byte, ttl time.Duration) error
Get(cf CFName, key []byte) ([]byte, bool, error)   // bool = found
Delete(cf CFName, key []byte) error
```

- **Point lookups only.** No range scans or iteration API, in v1. If the LLM
  cache use case later needs "list all keys for agent X", that's solved with
  a key-namespacing convention (`agentID:hash`).
- **Keys and values are raw `[]byte`.** Callers (the cache proxy) own serialization. Keeps the engine simple and reusable.
- **No transactions, no multi-key atomicity, no replication.** Single-key
  operations are atomic; that's the only guarantee. Explicitly out of scope.

## Durability decision

Every `Put` / `Delete` is logged to the WAL and **fsync'd as a batch of N operations**. Distributed system is a non-goal for this project plan, but will be the explicit focus in the next project plan.

## Column Families

Every method now takes a `CFName` and belongs to a column family.

`Reasoning:` The target use case (LLM cache) needs isolated keyspaces with different tuning -

- `exact cache` - small keys, high write rate and short TTL
- `stats` - tiny values, extremely high write rate, no durability urgency
- `metadata` - small values, rarely written with no expiry

_Design Choice_: Here `CFName` is a string instead of a handle.

## Delete semantics

A `Delete` is implemented as a **tombstone write** and not a physical removal.
This matters because:

- the memtable must distinguish "key absent" from "key explicitly deleted"
  once we add persistence, otherwise an older SSTable value could
  incorrectly "reappear" after a delete
- `Get` must treat a tombstone as "not found" even though an entry exists

## TTL semantics

TTL is stored as an absolute expiry timestamp (`time.Time`) alongside the
value. Expiry is checked lazily on read (`Get` returns not-found for an
expired key). A background sweeper is a compaction-time concern, revisited later.

## Explicit non-goals for the whole project

- No range scans / iterators
- No transactions or distributed semantics
- No replication or clustering
- No custom on-disk vector index (semantic search stays on pgvector)
