# Frontend implementation

The React/TypeScript/Vite frontend is implemented and deployed on
[Vercel](https://flowforge-weld-five.vercel.app), calling the
[Render API](https://flowforge-api-m3ua.onrender.com/health) directly.
See [frontend setup](../frontend/README.md) and [deployment](deployment.md).

## Pages and components

| Route | Source in `frontend/src/pages/` | Responsibility |
| --- | --- | --- |
| `/` | `PipelinesPage.tsx` | List, loading, empty/error states |
| `/pipelines/new` | `CreatePipelinePage.tsx` | Controlled pipeline/job form and create request |
| `/pipelines/:pipelineId` | `PipelinePage.tsx` | Definition, dependencies, Start Run |
| `/runs/:runId` | `RunPage.tsx` | Snapshot, job states, cancel/retry, live events |

`App.tsx` defines routes and redirects unknown paths to `/`. `AppShell`
provides shared navigation. `JobGraph` renders ordered job cards and dependency
labels; it is not a graphical node-and-edge editor. `EventLog` renders live events;
`Metric`, `StatusBadge`, and shared UI components provide presentation.

## REST flow

`src/api/types.ts` defines DTOs. `src/api/client.ts` owns JSON requests,
`ApiError` decoding, and `VITE_API_URL` (build-time; default localhost:8080).
`usePipelines` loads and refreshes the list; `usePipeline` loads one definition.

The creation form maps string fields into the request DTO and navigates to the
new pipeline. Durations use Go syntax; dependencies are comma-separated job IDs.
Backend validation is authoritative. Start Run navigates to the new run.
Cancel is available for active runs; Retry is available for failed/cancelled runs
and creates a new run without changing the original. Mutation errors are displayed.
There is no TanStack Query cache or authentication layer.

## SSE flow

`src/hooks/use-run.ts` implements `useRun`:

1. Fetch a current REST snapshot.
2. Open `EventSource` at the API's `/api/runs/<runId>/events`.
3. Replace the snapshot on `payload.run`; merge mutable state on `payload.job`.
4. Keep the newest 50 events for the visible list.
5. On error, close the stream, wait two seconds, fetch another snapshot, reconnect.
6. Close the stream and clear retry timers during cleanup.

The list resets on a new subscription lifecycle; it has no history API.
Opening a finished run may show completed jobs and an empty event panel.

The current hook does not deduplicate IDs or detect sequence gaps. There is a
snapshot/subscription race and no replay buffer; reload for an authoritative
snapshot if needed. Gap detection and replay are future hardening work.

## Hosting

Vercel uses root `frontend`. `vercel.json` supplies the direct-route SPA fallback.
The API is on Render; changing its hostname requires rebuilding the frontend.
Preview domains, other aliases, and local preview ports need explicit Render
CORS entries. Native `EventSource` follows cross-origin rules just like REST.

## Demo verification

[Import the 11-job fixture](examples/README.md) for dependency scheduling,
parallelism, timeouts, retry, and blocking. Only delay jobs run; fixture import
is manual. Verify creation, Start Run, events, cancel/retry, and direct URL reload.
After an API restart, old definition/run IDs may return `404`.
A `CONNECTED` stream does not imply that earlier events were replayed.
