# Iteration 1 - in memory memtable

## Decisions

- The Memtable only handles the lock for the column family map.
- Every data structure maintains its own mutex.
- Three available data structures for the memtable:
  - sorted slice
  - balanced trees
  - skiplist

## TODOs

- CF deletion and CF lifecycle
- Add and manage sequence number in Entry
