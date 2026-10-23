# FlowForge

FlowForge is a job orchestration demo with a Go API and a React/TypeScript
frontend. Define a directed acyclic graph of jobs, start independent runs, and
watch dependencies, parallel execution, retries, timeouts, and cancellation over
Server-Sent Events (SSE).

## Live application

| Component | Address | Hosting |
| --- | --- | --- |
| Frontend | [Open FlowForge](https://flowforge-weld-five.vercel.app) | Vercel, React/Vite |
| API | [Health endpoint](https://flowforge-api-m3ua.onrender.com/health) | Render, Go web service |

Deployment uses the `main` branch. See the [deployment guide](docs/deployment.md)
for exact settings, environment variables, release checks, and troubleshooting.

## Features and limits

- Create pipelines with DAG validation and a per-run concurrency limit.
- Start, cancel, and retry independent runs; inspect job states and attempts.
- View live execution events and recover REST snapshots after reconnecting.
- Import the [11-job demo](docs/examples/README.md) to explore failure and retries.

Only the in-process `delay` executor is implemented. Names such as "Build artifacts"
or "Deploy preview" simulate a CI/CD flow: they do not execute shell commands,
build a repository, or deploy another application.

This is a public MVP without authentication, a database, or durable event history.
Pipelines and runs live in one API process and disappear when it restarts or is
replaced. Demo definitions are in the repository but are **not loaded automatically**.
Import them again after a reset.

## Local development

Requirements: Go 1.22 or newer, Node.js 24 (matching the Vercel project), and pnpm
compatible with the version 9 lockfile.

Start the API in one terminal:

```sh
cd backend
go run ./cmd/api
```

Start the frontend in another terminal:

```sh
cd frontend
pnpm install --frozen-lockfile
pnpm dev
```

Open `http://localhost:5173`. The API defaults to `http://localhost:8080`;
no environment variables are needed for this local setup. If Vite chooses another
port, use port 5173 or add that exact origin to the API allowlist.

## Verification

From `backend/`:

```sh
go vet ./...
go test ./...
go test -race ./...
```

From `frontend/`:

```sh
pnpm build
pnpm lint
```

Render runs `go test ./...` before compiling the API. Race detection and
`go vet` are separate developer checks, not part of that build command.

## Documentation

- [Deployment and operations](docs/deployment.md)
- [Demo pipelines and import commands](docs/examples/README.md)
- [Frontend setup](frontend/README.md)
- [Frontend implementation and data flow](docs/frontend.md)
- [Architecture and ownership](docs/architecture.md)
- [Domain contracts](docs/domain.md)
- [HTTP API contract](docs/http-api.md)
- [Realtime events contract](docs/events.md)
- [Development status and next steps](docs/development-plan.md)
- [Future shell executor](docs/shell-executor.md)

## Repository structure

```text
backend/
  cmd/api/              composition, environment, HTTP server, shutdown
  internal/pipeline/    definitions, validation, in-memory repository
  internal/run/         snapshots, transitions, in-memory repository
  internal/scheduler/   concurrency, retries, timeouts, cancellation
  internal/runmanager/  lifecycle and active cancellation registry
  internal/events/      typed events and in-process subscriptions
  internal/httpapi/     REST/SSE handlers, DTOs, errors, CORS
frontend/
  src/api/              typed REST client and contracts
  src/hooks/            snapshots and SSE lifecycle
  src/pages/            list, creation, details, run view
  src/components/       layout, job cards, event stream, UI
  vercel.json           SPA fallback for direct application links
docs/
  examples/             JSON fixtures and import instructions
```
