# jQueue

A distributed job queue server written in Go, exposing a gRPC API for clients, workers, and a CLI.

Clients enqueue jobs (immediately or scheduled) and query their status. Workers lease jobs via long-polling, then Ack/Nack/ExtendLease using an opaque lease token. Delivery is at-least-once with no ordering guarantee; failed or timed-out jobs are retried with exponential backoff + jitter until a per-job retry limit, after which they land in a dead-letter queue (DLQ).

## Features

- Enqueue with client-supplied idempotency keys (duplicates are flagged, not re-created)
- Immediate or scheduled jobs
- Batch-capable leasing with long-polling (`wait_timeout`)
- Lease-based delivery: Ack / Nack / ExtendLease via lease token
- Automatic retries with exponential backoff + jitter, DLQ on retry exhaustion
- Queue administration: create, delete, purge, status, list
- Graceful shutdown (in-flight work drains before exit)
- Prometheus metrics endpoint (planned)

## Job lifecycle

```
Scheduled → Pending → Leased → Completed | Failed | DLQ | Cancelled
```

A periodic sweeper moves failed/timed-out jobs back to Pending or to the DLQ based on the retry limit. Completed jobs are cleaned up 24 hours after completion.

## gRPC services

| Service | Audience | RPCs |
|---|---|---|
| `JobService` | Clients | Enqueue, Status, Cancel, ListJobs, DLQRetry |
| `LeaseService` | Workers | LeaseJobs, Ack, Nack, ExtendLease |
| `QueueService` | Admin/CLI | Create, Delete, Purge, Status, List |
| `HealthService` | All | Liveness check |

## Getting started

```sh
# Build
go build ./...

# Run the server (listens on :50051)
go run .

# Run tests
go test ./...
```

gRPC reflection is enabled, so you can poke at the API with [grpcurl](https://github.com/fullstorydev/grpcurl):

```sh
grpcurl -plaintext localhost:50051 list
```

### Proto codegen

Protos live under `proto/jqueue/v1` and are managed with [buf](https://buf.build) (v2, managed mode). Generated Go code is committed under `gen/go/`.

```sh
buf lint      # lint protos
buf generate  # regenerate (wipes gen/go — never hand-edit generated files)
```

## Project layout

```
main.go            # entry point: listener, signal handling, graceful shutdown
internal/server/   # gRPC service implementations (one file per service)
internal/store/    # storage interface + in-memory implementation
proto/jqueue/v1/   # proto sources
gen/go/jqueue/v1/  # generated code (committed)
docs/plan.md       # full requirements and design
```

## Roadmap

1. **Phase 1 (current):** all state in memory, single server, no persistence
2. **Phase 2:** persistence via a shared SQL database, still a single broker
3. **Phase 3:** multiple distributed brokers over the shared database

The storage layer is behind an interface so the in-memory store can be swapped for SQL without touching the gRPC service layer.
