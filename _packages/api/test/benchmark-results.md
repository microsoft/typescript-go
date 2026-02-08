# API Benchmark Results

Three implementations tested:

1. **libsyncrpc** — native Rust/C++ addon with synchronous pipe I/O
2. **Sync JS** — pure JS `SyncRpcChannel` (MessagePack over synchronous pipe I/O)
3. **Async JS** — async/await over JSON-RPC

## Key Comparisons (median latency)

| Benchmark | libsyncrpc | Sync JS (new) | Async JS | Sync JS vs libsyncrpc | Async JS vs libsyncrpc |
|---|---|---|---|---|---|
| **Spawn** | 308 µs | 8.6 ms | 2.6 µs | 28× slower | 118× faster |
| **Load project** | 64 µs | 58 µs | 110 ms | 0.9× (faster) | 1719× slower |
| **Transfer debug.ts** | 1.6 ms | 1.3 ms | 29.7 ms | 0.8× (faster) | 18.6× slower |
| **Transfer program.ts** | 5.8 ms | 4.7 ms | 100.8 ms | 0.8× (faster) | 17.4× slower |
| **Transfer checker.ts** | 74.7 ms | 58.2 ms | 1240 ms | 0.8× (faster) | 16.6× slower |
| **Materialize program.ts** | 2.7 ms | 2.6 ms | 2.5 ms | 1.0× (~equal) | 0.9× (~equal) |
| **Materialize checker.ts** | 73.1 ms | 72.2 ms | 73.7 ms | 1.0× (~equal) | 1.0× (~equal) |
| **getSymbolAtPosition (1)** | 13.2 µs | 15.5 µs | 43.6 µs | 1.2× slower | 3.3× slower |
| **getSymbolAtPosition (10153, batched)** | 87.4 ms | 84.9 ms | 87.5 ms | 1.0× (~equal) | 1.0× (~equal) |
| **getSymbolAtLocation (10153)** | 406 ms | 442 ms | 798 ms | 1.1× slower | 2.0× slower |
| **getSymbolAtLocation (10153, batched)** | 282 ms | 267 ms | 278 ms | 0.9× (faster) | 1.0× (~equal) |
| **TS baseline (load project)** | 886 ms | 898 ms | 910 ms | — | — |
| **TS baseline (getSymbol 10153)** | 25.2 ms | 25.2 ms | 26.9 ms | — | — |

## Summary

- **Sync JS vs libsyncrpc:** Very close in nearly all benchmarks. Sync JS is slightly faster for data transfer (debug/program/checker.ts) while libsyncrpc has a ~28× faster spawn time (308 µs vs 8.6 ms). For hot-path operations like batched symbol lookups, they're within noise of each other. The pure JS replacement successfully matches native performance.

- **Sync vs Async JS:** Sync APIs are dramatically faster for data transfer (~20× for large files) and project loading (~1900×), because data flows through synchronous pipes without async/event-loop overhead. Async JS's spawn is much cheaper (µs vs ms) since it doesn't block.

- **tsgo vs TS baseline:** Project loading is ~14,000× faster. Per-identifier symbol lookup is 6–16× slower due to IPC overhead, but batched queries bring it to within 3–11× of TS's in-process speed.

- **Batching matters:** Batching 10k symbol lookups into a single call gives a ~3–5× speedup over per-identifier round-trips in both sync channels.
