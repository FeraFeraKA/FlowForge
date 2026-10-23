# FlowForge

FlowForge is a small job orchestration engine written in Go. A pipeline describes
a directed acyclic graph of jobs; every pipeline start creates an independent run
whose jobs are scheduled according to dependencies, retries, timeouts,
cancellation, and a concurrency limit.

The backend is the only implementation target during the MVP. The future frontend
will use REST for commands and snapshots and Server-Sent Events (SSE) for live
updates.

## Current stage

The backend MVP is implemented: DAG validation, in-memory repositories, explicit
state transitions, sequential and concurrent schedulers, retries, timeouts,
cancellation, REST, SSE, local CORS, and graceful shutdown. The test suite passes
with Go's race detector.

Run it locally:

```bash
cd backend
go run ./cmd/api
```

The API listens on `http://localhost:8080`. Request examples and frontend-facing
contracts live in the documents below.

## Documents

## Deployment

The React/Vite frontend is deployed on Vercel with root directory `frontend`.
Set `VITE_API_URL` to the public Render API origin before building. The
`frontend/vercel.json` rewrite supports direct links to application pages.

The Go API runs on Render with these commands from the repository root:

```sh
cd backend && go test ./... && go build -o bin/api ./cmd/api
```

Start command: `cd backend && ./bin/api`. Render provides `PORT`; local runs
default to port 8080. Set `ALLOWED_ORIGINS` to a comma-separated list of allowed
frontend origins, for example `https://your-project.vercel.app`. Without this
variable, only the local Vite development origins are allowed. `GET /health`
returns the API health status.

Repositories are in memory: pipelines and runs are lost on API restarts,
redeployments, and free-service spin-downs. This deployment retains the current
MVP storage model and does not provision a database.

- [Architecture and ownership](docs/architecture.md)
- [Domain contracts](docs/domain.md)
- [HTTP API contract](docs/http-api.md)
- [Realtime events contract](docs/events.md)
- [Frontend handoff](docs/frontend.md)
- [Development plan](docs/development-plan.md)
- [Post-MVP shell executor](docs/shell-executor.md)

## Backend structure

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
