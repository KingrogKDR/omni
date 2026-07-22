# Goals and Design for Omni

## Goals:

- Simple architecture
- Embedded datastore
- Human-readable tooling
- Zero/optional configuration
- Predictable performance
- Bounded recovery time
- Minimal storage cost for only the latest data

## Storage Engine design:

- Two layers:
  - upper layer:
    - durable storage
    - append-only
    - stores the latest versions that have not yet been compacted into the lower layer.
    - contains segments of fixed size (size is configurable, by default: 64MB)
    - a segment should only be deleted after its content is successfully compacted into the lower layer.
  - lower layer:
    - permanent compact snapshot of the latest sync
    - only contains the latest synced data
    - it will be stored in the form of partitions with hash(key)%n where n -> no. of partitions

# Read/Write design:

- Read:
  - read always occurs top-down (upper -> lower)
  - uses paritions for better read performance because we only search the partition for the key asked.
  - tombstones are never read
- Write:
  - write always occurs in the upper layer.
  - it never touches the lower layer.
  - only the sync process writes from the upper layer to the lower layer in a compact format
  - deleting marks the key a tombstone (FLAG) in both the upper layer and the in-mem map
  - duplication is never allowed after sync

# Core Principles:

- Foreground writes are always sequential.
- Acknowledged/Synced writes are durable.
- Background/Synchronization work never blocks writes.
- Only the latest live version of a key occupies long-term storage.
- The architecture should be explainable in one diagram.
