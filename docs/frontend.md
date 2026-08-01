# Frontend handoff

This document separates beginner-friendly React work from integration work. The
frontend is planned for React, TypeScript, Vite, and `pnpm`. It has no dependency
on a database or authentication.

## Suggested application pages

1. `/pipelines` — pipeline list, empty/error/loading states, and create action.
2. `/pipelines/new` — controlled form for pipeline and jobs.
3. `/pipelines/:pipelineId` — pipeline details and start-run action.
4. `/runs/:runId` — live run status, jobs, dependencies, attempts, and logs.

## Work suitable for the second frontend developer

These components can initially be built with typed mock props and Storybook or a
small local mock-data module; they do not need networking knowledge.

- `AppShell` and page layout;
- `StatusBadge` for every documented run/job status;
- `ProgressBar` with accessible label and 0–100 clamping;
- `PipelineCard` and `PipelineList`;
- loading, empty, and error-state components;
- `PipelineForm` with controlled name and `maxParallel` inputs;
- `JobDefinitionForm` for ID, name, type, duration, timeout, retries, retry delay;
- add/remove/reorder job form rows;
- dependency multi-select that prevents obvious self-selection;
- client-side field-error display using API `details` paths;
- pipeline detail summary and dependency list;
- `RunHeader`, job cards, attempt counter, timestamps, and error display;
- log list and basic filters, using mock log entries until backend logging exists;
- confirmation modal for Run cancellation;
- responsive and accessible keyboard/focus behavior.

The form may provide friendly client validation, but the backend remains the
authority for graph validation and executor configuration.

## Integration work

The primary developer can own:

- typed REST client and API error decoding;
- optional TanStack Query setup for REST snapshots and mutation invalidation;
- pipeline creation and start/cancel/retry mutations;
- `useRunEvents(runId)` SSE lifecycle;
- a reducer that replaces state from REST and applies typed events;
- reconnection with backoff and snapshot refresh;
- detection of duplicate/out-of-order event IDs and event-ID gaps;
- synchronization between query cache and realtime state;
- rendering the DAG/dependency relationships;
- complex Run-page state and terminal-state handling.

## Recommended TypeScript contracts

Define explicit frontend types matching the documented DTOs, including string
unions for statuses and event types. Do not share generated Go structs directly.
Keep three layers where useful:

- API DTOs matching wire JSON;
- form state, where duration fields may temporarily be empty/invalid strings;
- view models derived for display.

## Mock-first workflow

Create fixtures for at least:

- an empty pipeline list;
- a valid diamond-shaped pipeline;
- pending and running Runs;
- a successful Run;
- a failed Run with blocked downstream jobs;
- a retrying job on attempt 2;
- a cancelled Run;
- validation and conflict error envelopes.

Mock event sequences should include success, retry-then-success, final failure with
blocked jobs, cancellation, duplicate event, and reconnect snapshot replacement.

## Suggested reducer actions

The reducer does not need to mirror raw event names exactly. A small stable set is
enough:

```text
snapshotReceived
eventReceived
connectionChanged
streamError
```

`snapshotReceived` replaces authoritative Run state. `eventReceived` replaces the
affected JobRun or Run summary using the complete entity carried by the payload.

## Local development assumptions

- frontend default origin: `http://localhost:5173`;
- backend base URL is configuration, not hardcoded throughout components;
- all commands use `pnpm`, never npm;
- `EventSource` cannot set arbitrary headers in the usual browser API, which is
  acceptable because the MVP has no authentication;
- UI must not assume SSE delivers historical events.

## Definition of frontend readiness

The backend is ready for integration when it provides stable examples for every
endpoint, a complete Run snapshot, the event envelope and payloads, CORS, predictable
error codes, and a scripted demo pipeline covering success/failure/retry/cancel.
