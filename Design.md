# Design document for Omni

Everything related to the design decisions made keeping in mind the tradeoffs and usecases.

## The network layer

- Wire format: Protobuf

#### Context

Omni nodes exchange messages over the network. These nodes won't always run identical binary versions at the same time. The format needs version mismatches to fail immediately at the boundary and not pass silently, corrupting application logic.

- Options considered:
  - Plain JSON — no schema, fields are whatever keys happen to be present
  - Typed JSON — JSON Schema / struct validation layered on top
  - Protobuf — schema-defined binary format

Plain JSON does no enforcement of types which is a problem in itself for our usecase
Typed JSON moves the typing outside the wire format which introduces "dependencies" and needs tracking and versioning of these dependencies, leading to high maintenance requirements.
Protobuf solves all these problems by enforcing type safety into the wire format itself. A type violation fails at deserialization, at the network boundary, not inside application logic.

Benchmarked on a representative Raft command message:

| Operation             |       JSON |  Protobuf |       Improvement |
| --------------------- | ---------: | --------: | ----------------: |
| Marshal               |  418 ns/op | 182 ns/op |  **2.30× faster** |
| Unmarshal             | 2171 ns/op | 208 ns/op | **10.44× faster** |
| Marshal memory        |      192 B |      96 B |    **2.00× less** |
| Unmarshal memory      |      400 B |     216 B |    **1.85× less** |
| Unmarshal allocations |          7 |         3 |   **2.33× fewer** |
| Wire size             |      190 B |      89 B | **2.13× smaller** |

Still there do exist some tradeoffs to be aware of:

- Codegen step is an extra moving part in the build which JSON didn't need
- Binary format — not human-readable; debugging needs `protoc/buf` to inspect a message
- Schema changes require regenerating code across every consumer
- Proto3 fields don't track presence by default — an unset field and a field explicitly set to its zero value decode identically. Every field where "was this explicitly set" matters (increasingly relevant once versioning/MVCC enters the picture). Needs the optional keyword to opt into presence tracking.
- Benchmark numbers above are single-machine, single-message-shape — not a substitute for measuring under real concurrent load

Yet Protobuf wins in both safety and speed for Omni.

Full reasoning in the [linked post](https://kingrogkdr.github.io/post.html?slug=omni-post-1)
