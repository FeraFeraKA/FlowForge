# Development plan

## Current deployed state

The Go API and React/Vite frontend are implemented and deployed on Render and
Vercel. See [deployment and operations](deployment.md) for service settings,
environment variables, release checks, and data-loss behavior.

The stages below record the backend implementation. They are not a list of
missing work. Current frontend behavior is described in [frontend](frontend.md).

Each stage starts by stating its goal, files, responsibilities, data flow,
goroutines, state owner, error handling, verification commands, and one manual
test. A stage is reviewed before the next one begins.

## Stage 0 — contracts and skeleton

Status: complete. Contracts and package boundaries were established before
implementation.

## Stage 1 — process skeleton and health

Status: complete.

Create the Go module, `cmd/api/main.go`, the standard `http.ServeMux`, `/health`,
application context, and bounded graceful shutdown. Verify with `go run`, `curl`,
`go test`, `go vet`, and a manual SIGINT test.

## Stage 2 — Pipeline model and DAG validation

Status: complete.

Implement data types and pure validation with table-driven tests for duplicates,
unknown dependencies, self-dependencies, cycles, and valid graphs.

## Stage 3 — Pipeline repository and service

Status: complete.

Implement the mutex-protected in-memory repository and prove through tests that
saved and returned values cannot mutate repository-owned data.

## Stage 4 — Run model and state machine

Status: complete.

Implement explicit legal transitions and table-driven tests. There are no
goroutines yet.

## Stage 5 — DelayExecutor

Status: complete.

Implement timer completion and context cancellation. Test success, cancellation,
and timeout deterministically.

## Stage 6 — sequential scheduler

Status: complete.

Implement dependency readiness, retries, blocking, cancellation, and terminal Run
status with one job executing at a time. Use deterministic fake executors in tests.

## Stage 7 — Run repository and run manager

Status: complete.

Add deep snapshots, application-derived run contexts, active cancel registry, and
cleanup after terminal completion.

## Stage 8 — concurrent scheduler

Status: complete.

Keep one coordinator as state owner. Add bounded concurrent attempts, result
channels, `MaxParallel`, explicit completion, cancellation, and race tests.

## Stage 9 — event bus

Status: complete.

Add event sequencing, bounded subscriber channels, non-blocking publication, and
safe subscribe/unsubscribe ownership.

## Stage 10 — REST API

Status: complete.

Implement DTO mapping, error envelope, pipeline/run commands and queries, and
explicit local-development CORS.

Deployment now supplies allowed frontend origins through `ALLOWED_ORIGINS`.

## Stage 11 — SSE

Status: complete.

Add streaming, flush, heartbeat, disconnect cleanup, and manual `curl -N` tests.

## Stage 12 — integration and hardening

Status: complete.

Test success, parallelism, retry, timeout, blocking, cancel, reconnect, and
shutdown. Run:

```bash
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
```

Document final examples and provide frontend demo fixtures.

## Stage 13 — frontend and hosting

Status: implemented and deployed.

The frontend lists/creates pipelines, shows definitions, starts runs, and displays
live job states/events with cancel/retry actions. Vercel hosts the static bundle
with direct-route rewrites. Render builds/tests the Go API and starts it using
`PORT`; `ALLOWED_ORIGINS` authorizes the frontend origin.

Deployment verification covered health, creation, successful execution, SSE,
and page reload on a direct route. The [complex fixture](examples/README.md)
demonstrates timeout/retry/blocking and is imported manually.

## Remaining work

- Durable storage for pipelines and run snapshots; restart currently clears them.
- Optional idempotent demo seeding; no startup import exists today.
- Authentication/authorization before using the public demo for private workloads.
- SSE ID deduplication, gap detection, and replay.
- A graphical DAG editor if needed; the UI currently renders job/dependency cards.
- Durable job logs and artifacts; the event panel is a live orchestration feed.
- External monitoring and a resolved access path for Vercel runtime logs.
- Isolated command execution only after the [shell-executor design](shell-executor.md).
