# Architecture and ownership

## Boundaries

FlowForge deploys a static React/Vite frontend on Vercel and one Go API process
on Render. The browser sends REST and SSE requests directly to Render over HTTPS.
The API is not a distributed scheduler and has no persistence across restarts.
See [deployment and operations](deployment.md) for live endpoints and settings.

The implemented backend packages are:

- `pipeline`: saved pipeline definitions and DAG validation;
- `run`: run/job-run data and legal state transitions;
- `scheduler`: readiness decisions, retries, execution, and terminal-state logic;
- `runmanager`: starts long-lived runs and owns their cancellation functions;
- `events`: typed events and subscriber lifecycle;
- `httpapi`: request/response DTOs, REST handlers, SSE, errors, and CORS;
- `cmd/api`: dependency construction, signals, HTTP server, and graceful shutdown.

Dependencies should point inward from process/transport code. Domain packages do
not import HTTP packages. The scheduler does not know about handlers or SSE.

## Pipeline immutability

Pipelines cannot be edited in the MVP. A changed definition is created as a new
pipeline. Starting a run copies the effective job definitions into the run
snapshot, so an execution is self-contained and reproducible.

## Active Run ownership

An active run has exactly one writer: its scheduler coordinator.

- Workers receive immutable attempt input and return an execution result.
- Workers never mutate `Run` or `JobRun`.
- The coordinator validates and applies state transitions serially.
- After mutation, the coordinator stores a deep snapshot and publishes the
  corresponding event, in that order.
- Repositories never expose internal maps, slices, or mutable pointers. `Save`
  stores a deep copy; `Get` and `List` return deep copies.
- HTTP handlers and SSE code cannot mutate domain state.

The sequential scheduler follows the same ownership rule without worker
goroutines. The deployed API uses the concurrent scheduler with the delay executor.

Repository mutexes protect repository maps and copying operations. They are not a
substitute for ownership of an object after a lock is released.

## Context ownership

There are three different lifetimes:

1. The HTTP request context owns only the start command.
2. A run context owns the long-lived execution.
3. The application context owns every active run.

`StartRun(ctx, pipelineID)` uses the request context while validating the command,
loading the pipeline, creating the initial run, and registering it. The background
execution does not inherit that request context.

The run manager creates each run context from the application context. It stores
the matching cancel function in a private active-run registry. A cancel request
invokes that function. A job attempt with a timeout derives a timeout context from
the run context. Retry waits also observe the run context.

When a run reaches a terminal state, the run manager removes its cancel function
and signals its completion. SIGINT or SIGTERM cancels the application context and
therefore its run contexts. The HTTP server gets a ten-second graceful-shutdown
deadline; deferred cleanup waits for the run manager and closes the event bus.
That ten-second deadline applies to HTTP shutdown, not a separate timeout around
the run-manager wait. Delay execution and retry waits honor cancellation.

Business logic must not replace an incoming context with `context.Background()`.

## HTTP and realtime consistency

REST snapshots are the source of truth. SSE provides low-latency notification but
is not durable storage.

For a state change, the coordinator performs:

1. validate the transition;
2. mutate its private state;
3. save a deep snapshot;
4. publish an event describing the committed state.

The MVP does not promise transactional atomicity between the repository and event
bus. A reconnecting frontend fetches a fresh REST snapshot before resubscribing.

## Router

The backend targets Go 1.22 or newer and uses the standard `http.ServeMux` with
method-aware patterns and `Request.PathValue`. Chi is not needed for the current
route set. Middleware remains ordinary `func(http.Handler) http.Handler`.

## Hosting boundary

The API listens on all interfaces at the port supplied by `PORT`, defaulting
to 8080 locally. `ALLOWED_ORIGINS` configures exact frontend origins.
`VITE_API_URL` is compiled into the frontend bundle. Vercel's SPA rewrite
handles application routes, while Render owns API routing and SSE connections.

There is no database, external queue, authentication, startup demo seeding, or
shell execution. Definitions, snapshots, active cancellation state, and events
belong to the current API process. Multiple API replicas would not share them.
Redeploying or restarting the process clears the demo state.
