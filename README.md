# FlowForge

FlowForge is a small job orchestration engine written in Go. A pipeline describes
a directed acyclic graph of jobs; every pipeline start creates an independent run
whose jobs are scheduled according to dependencies, retries, timeouts,
cancellation, and a concurrency limit.

The backend is the only implementation target during the MVP. The future frontend
will use REST for commands and snapshots and Server-Sent Events (SSE) for live
updates.

## Current stage

Stage 0 contains contracts and the planned package structure only. There is no Go
implementation yet.

## Documents

- [Architecture and ownership](docs/architecture.md)
- [Domain contracts](docs/domain.md)
- [HTTP API contract](docs/http-api.md)
- [Realtime events contract](docs/events.md)
- [Frontend handoff](docs/frontend.md)
- [Development plan](docs/development-plan.md)
- [Post-MVP shell executor](docs/shell-executor.md)

## Planned backend structure

```text
backend/
  cmd/api/              application composition and process lifecycle
  internal/pipeline/    pipeline model, DAG validation, repository, service
  internal/run/         run model, snapshots, and state transitions
  internal/scheduler/   scheduling algorithm and executors
  internal/runmanager/  long-lived run contexts and active-run registry
  internal/events/      event model and in-process subscriptions
  internal/httpapi/     REST, SSE, DTOs, errors, middleware, CORS
```

The directories exist as placeholders. Files will be introduced one completed
stage at a time.
