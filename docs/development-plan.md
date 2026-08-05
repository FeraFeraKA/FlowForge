# Development plan

Each stage starts by stating its goal, files, responsibilities, data flow,
goroutines, state owner, error handling, verification commands, and one manual
test. A stage is reviewed before the next one begins.

## Stage 0 — contracts and skeleton

Status: complete when these documents and empty package directories are reviewed.
No Go implementation or dependency is introduced.

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

Add event sequencing, bounded subscriber channels, non-blocking publication, and
safe subscribe/unsubscribe ownership.

## Stage 10 — REST API

Implement DTO mapping, error envelope, pipeline/run commands and queries, and
explicit local-development CORS.

## Stage 11 — SSE

Add streaming, flush, heartbeat, disconnect cleanup, and manual `curl -N` tests.

## Stage 12 — integration and hardening

Test success, parallelism, retry, timeout, blocking, cancel, reconnect, and
shutdown. Run:

```bash
go fmt ./...
go vet ./...
go test ./...
go test -race ./...
```

Document final examples and provide frontend demo fixtures.
