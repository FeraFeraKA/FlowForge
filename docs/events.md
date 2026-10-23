# Realtime event contract

The deployed stream is on the Render API:
`https://flowforge-api-m3ua.onrender.com/api/runs/{runID}/events`.
Use a current run ID returned by the API. The frontend connects directly from
Vercel; see [deployment](deployment.md) for origins and [frontend](frontend.md)
for the implemented reconnect behavior.

## SSE endpoint

`GET /api/runs/{runID}/events` returns:

```text
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
```

Each message uses an SSE `id`, `event`, and JSON `data` field:

```text
id: 42
event: job.started
data: {"id":"42","type":"job.started","runId":"run_opaque","jobId":"tests","timestamp":"2026-08-01T12:01:02Z","payload":{"job":{"jobId":"tests","status":"running","attempt":1}}}

```

## Envelope

```json
{
  "id": "42",
  "type": "job.started",
  "runId": "run_opaque",
  "jobId": "tests",
  "timestamp": "2026-08-01T12:01:02Z",
  "payload": {}
}
```

- `id` is monotonically increasing within one run, not globally;
- `jobId` is omitted for run-level events;
- `type` is duplicated in the SSE `event` field so `EventSource` listeners can
  subscribe by type;
- run payloads contain the complete REST-style run snapshot; job payloads contain
  `jobId`, `status`, `attempt`, `progress`, `error`, `startedAt`, and
  `finishedAt`. Job definition fields remain in the REST snapshot;
- clients ignore unknown event types and unknown payload fields.

Initial event types:

- `run.started`, `run.completed`, `run.failed`, `run.cancelled`;
- `job.queued`, `job.started`, `job.progress`, `job.retrying`;
- `job.completed`, `job.failed`, `job.blocked`, `job.cancelled`;
- `job.log` is reserved for the later logging implementation.

## Subscription lifecycle

- A subscription is scoped to one Run.
- The request context owns the SSE handler and unsubscribe operation.
- The event bus owns each subscriber channel and is the only component that
  closes it.
- The SSE handler never closes the channel.
- Disconnect, write failure, request cancellation, or server shutdown removes the
  subscription.
- The handler writes an SSE comment heartbeat every 15 seconds. This interval is
  currently a code constant, not an environment setting.

## Slow clients

Each subscriber has a bounded channel with an initial capacity of 64 events.
Publishing never blocks the scheduler. If a subscriber buffer is full, the bus
unsubscribes and closes that subscriber. The client reconnects and restores state
from REST.

## Reconnection and replay

The MVP event bus is in memory and does not provide replay, including replay based
on `Last-Event-ID`. Event IDs are still present for ordering and future evolution.

The current frontend reconnects as follows:

1. close or observe loss of the old stream;
2. fetch `GET /api/runs/{runID}`;
3. replace local state with that authoritative snapshot;
4. open a new SSE stream;
5. apply newly received events as delivered and retain the latest 50 in the UI.

It waits two seconds before reconnecting after a stream error. ID deduplication
and gap detection are not implemented in the current hook.

There is a small snapshot/subscription race in this MVP: an event can occur after
the snapshot and before subscription. The next received event carries the complete
affected job's mutable state. Reconnection refetches the run; a manual page reload
also restores a snapshot if the UI becomes stale. A later version may add gap
detection, a bounded per-run replay buffer, and
`Last-Event-ID` support.

## Source of truth

SSE is not the only representation of state. Opening a Run page starts from a
REST snapshot, but that snapshot is also in memory: neither layer survives an API
process restart. Events improve freshness but do not replace snapshots. Opening
a finished run may show an empty event log because past events are not replayed.

To inspect a live run from a terminal (replace RUN_ID):

```sh
curl -N https://flowforge-api-m3ua.onrender.com/api/runs/RUN_ID/events
```

Connect while the run is active to observe transitions. A connected stream for
a finished run can remain open with heartbeats and no new execution events.
