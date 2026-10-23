# HTTP API contract

## Base URLs

- Deployed API: `https://flowforge-api-m3ua.onrender.com`.
- Local API: `http://localhost:8080`, unless `PORT` changes it.
- Deployed frontend: `https://flowforge-weld-five.vercel.app`.

All paths below are relative to the API origin, not the frontend origin.
The API currently has no authentication or cookies. CORS is a browser-origin
policy, not access control. Data is in memory; IDs may cease to exist after an
API restart. See [deployment](deployment.md) and [demo import](examples/README.md).

All request and response bodies use JSON and camelCase fields. Unknown fields in
command payloads are rejected. Unless noted otherwise, successful responses use
`Content-Type: application/json`.

## Error envelope

```json
{
  "error": {
    "code": "validation_failed",
    "message": "pipeline validation failed",
    "details": [
      {
        "field": "jobs[1].dependsOn[0]",
        "code": "unknown_dependency",
        "message": "dependency \"missing\" does not exist"
      }
    ]
  }
}
```

`details` is optional. Expected domain errors may expose safe validation details.
Unexpected internal errors use code `internal_error` and a generic message; stack
traces and wrapped technical errors are never returned to clients.

Status mapping:

- malformed JSON, wrong JSON types, trailing JSON values: `400`;
- well-formed but semantically invalid command: `422`;
- missing pipeline or run: `404`;
- illegal state transition or retry/cancel conflict: `409`;
- unsupported content type: `415`;
- unexpected error: `500`.

## Health

### `GET /health`

Response `200`:

```json
{"status":"ok"}
```

## Pipelines

### `POST /api/pipelines`

Request:

```json
{
  "name": "Build and deploy",
  "maxParallel": 2,
  "jobs": [
    {
      "id": "install",
      "name": "Install dependencies",
      "type": "delay",
      "dependsOn": [],
      "timeout": "3s",
      "maxRetries": 0,
      "retryDelay": "0s",
      "config": {"duration": "1s"}
    },
    {
      "id": "tests",
      "name": "Run tests",
      "type": "delay",
      "dependsOn": ["install"],
      "timeout": "5s",
      "maxRetries": 2,
      "retryDelay": "500ms",
      "config": {"duration": "2s"}
    }
  ]
}
```

Response `201`: complete Pipeline resource. Header `Location` points to
`/api/pipelines/{pipelineID}`.

### `GET /api/pipelines`

Response `200`:

```json
{
  "items": [],
  "total": 0
}
```

MVP has no pagination. Items are sorted newest first, with ID as a deterministic
tie-breaker.

### `GET /api/pipelines/{pipelineID}`

Response `200`: complete Pipeline resource. Missing pipeline: `404`.

## Runs

### `POST /api/pipelines/{pipelineID}/runs`

The request has no JSON body. Response `202`: initial complete Run snapshot.
Header `Location` points to `/api/runs/{runID}`. The run continues independently
of the request connection.

### `GET /api/runs/{runID}`

Response `200`: current complete Run snapshot, including its JobRuns and effective
job definitions. Missing run: `404`.

### `POST /api/runs/{runID}/cancel`

The request has no JSON body. Response `202`: the latest Run snapshot after the
cancellation request has been accepted. Cancellation is asynchronous; clients use
REST/SSE to observe the terminal snapshot.

Returns `409` if the run is already terminal. Repeated cancellation while the same
run is still cancelling is idempotently accepted with `202`.

### `POST /api/runs/{runID}/retry`

The request has no JSON body. This creates a new run from the old run's stored
definition snapshot. The original run remains unchanged.

Allowed only when the source run is `failed` or `cancelled`; otherwise returns
`409`. Response `202`: the new Run snapshot. `Location` points to the new run.

### `GET /api/runs/{runID}/events`

Opens an SSE stream as specified in `docs/events.md`. Missing run: `404`.

## CORS

The process reads `ALLOWED_ORIGINS` at startup as a comma-separated list of exact
origins. Without it, defaults are `http://localhost:5173` and
`http://127.0.0.1:5173`. A non-empty value replaces those defaults.
The deployed allowlist contains `https://flowforge-weld-five.vercel.app`.

The middleware allows `GET`, `POST`, and `OPTIONS`, plus `Content-Type` and
`Last-Event-ID`. A request with an unlisted `Origin` returns `403` with code
`origin_not_allowed`; requests without `Origin` proceed normally.
Allowed preflight requests return `204`. Wildcard origins and credentials are
not used. Every preview URL or additional frontend alias needs an explicit entry.
The same middleware protects the SSE route.

## DTO boundary

HTTP request DTOs are decoded separately from domain models. Responses are also
mapped explicitly. Domain structs therefore do not define the public API merely
because their fields happen to serialize to JSON.
