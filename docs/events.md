# Realtime event contract

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
- payloads contain the current affected `run` summary or complete affected
  `job`, not an ambiguous JSON patch;
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
- The handler periodically writes an SSE comment heartbeat so dead connections can
  be discovered. The exact interval will be configuration, with 15 seconds as the
  initial default.

## Slow clients

Each subscriber has a bounded channel with an initial capacity of 64 events.
Publishing never blocks the scheduler. If a subscriber buffer is full, the bus
unsubscribes and closes that subscriber. The client reconnects and restores state
from REST.

## Reconnection and replay

The MVP event bus is in memory and does not provide replay, including replay based
on `Last-Event-ID`. Event IDs are still present for ordering and future evolution.

Frontend reconnection sequence:

1. close or observe loss of the old stream;
2. fetch `GET /api/runs/{runID}`;
3. replace local state with that authoritative snapshot;
4. open a new SSE stream;
5. apply subsequently received events in increasing event-ID order.

There is a small snapshot/subscription race in this MVP: an event can occur after
the snapshot and before subscription. The next received event carries the complete
affected JobRun, and terminal UI should refetch after reconnect or when it detects
an event-ID gap. A later version may add a bounded per-run replay buffer and
`Last-Event-ID` support.

## Source of truth

SSE is never the only persisted representation of state. Opening a Run page always
starts from a REST snapshot. Events improve freshness but do not replace it.
