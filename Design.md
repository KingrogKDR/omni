# Design Doc for Omni

Omni is an embedded key-value storage engine designed around one primary goal:

> Provide a simple, predictable, and durable storage engine whose architecture is easy to understand while maintaining high write throughput, fast point reads, bounded crash recovery, and minimal long-term storage amplification.

Omni is not an LSM implementation nor an in-memory cache. Instead, it is a log-structured embedded KV store built around a two-layer storage model.

The design intentionally favors architectural simplicity over implementing every optimization found in mature databases.

# Design Philosophy

Modern storage engines often trade simplicity for configurability and maximum throughput.

Omni takes a different approach.

The core philosophy is:

- Sequential writes should always be cheap.
- Durability should never depend on background work.
- Background work should never block foreground writes.
- Long-term storage should contain only the latest live version of every key.
- Recovery should always be deterministic and bounded.
- Every major subsystem should have one clearly defined responsibility.

If a feature increases complexity without significantly improving these principles, it does not belong in Omni.

## Goals

- Embedded datastore
- Simple architecture
- Predictable performance
- Human-readable tooling
- Zero mandatory configuration
- Bounded crash recovery
- Durable writes
- Low storage amplification
- Future-ready for distributed replication

## Non Goals

Version 1 intentionally does not target:

- SQL support
- Transactions
- MVCC
- Rollback/version history
- Multi-node replication
- Secondary indexes
- Analytical workloads

These may be explored in future versions but are not part of the initial architecture.

## Storage Model:

Omni separates storage into two independent layers.

                 Client
                    │
                    ▼
         Append-only Upper Layer
                    │
          Background Synchronization
                    │
                    ▼
      Compact Partitioned Lower Layer

Each layer has a single responsibility.

- Upper layer:

The upper layer is the write-optimized portion of the database.

It accepts every mutation.

Responsibilites:

- Accept all writes
- Store recent updates
- Store tombstones
- Organize data into immutable segments of fixed size (size is configurable, by default: 64MB)

Properties:

- Durable
- Append-only
- Segmented
- Immutable after sealing

Once a segment becomes full it is sealed.

Sealed segments are never modified again.

- Lower layer:

The lower layer represents the compact, long-term view of the database.

Unlike the upper layer, it contains only the newest synchronized version of every live key.

Older versions never exist here.

Responsibilities

- Compact storage
- Fast reads
- Minimal storage amplification
- Reduce recovery time
- Partitioning

The lower layer is divided into partitions.

`partition = hash(key) % N` -> N = no.of partitions

Each partition stores only keys belonging to that partition.

Advantages:

- Localized synchronization
- Smaller rewrite operations
- Better cache locality
- Faster point lookups
- Easier parallelization in future

## Segment Lifecycle

Every upper-layer segment follows a deterministic lifecycle.

ACTIVE
│
▼
SEALED
│
▼
QUEUED_FOR_SYNC
│
▼
SYNCING
│
▼
SYNCED
│
▼
READY_FOR_DELETE
│
▼
DELETED

Each state transition is atomic.

Recovery only needs to inspect segments that have not reached the DELETED state

## Read/Write design:

- Read Path:

Reads always follow a merged-view model.

```
        GET(key)
           │
           ▼
    Check Upper Layer
           │
        Found?
        /    \
      Yes     No
      │        │
    Return   Check Lower
                │
                ▼
            Return
```

The upper layer always overrides the lower layer.

This guarantees read-your-own-write consistency.

- Write Path:

Every mutation is appended to the active upper-layer segment.

```
SET

↓

Append Record

↓

fsync

↓

ACK
```

Foreground writes never modify the lower layer.

This keeps write latency predictable.

# Core Principles:

- Foreground writes are always sequential.
- Acknowledged/Synced writes are durable.
- Background/Synchronization work never blocks writes.
- Only the latest live version of a key occupies long-term storage.
- The architecture should be explainable in one diagram.
