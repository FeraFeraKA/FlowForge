# Demo pipelines

The repository includes [complex-pipeline.json](complex-pipeline.json), an
11-job orchestration demo. It is **not imported automatically** at startup.
The live app may contain manually created demos, but their generated IDs and runs
disappear when the API process restarts. These instructions recreate the demo.

## What it demonstrates

- Two root jobs: Fetch source and Restore cache.
- Install dependencies waits for both roots.
- Lint, Unit tests, and Build artifacts fan out, with at most three active
  attempts per run.
- Successful branches continue through security scanning, packaging, integration,
  and publishing a report.
- Unit tests intentionally runs for four seconds with a one-second timeout.
  Two retries permit three total attempts, each separated by a 750ms retry delay.
- After the final failure, Deploy preview becomes blocked.
- The run finishes as `failed` once independent work has settled.

All jobs use the delay executor. These names simulate CI/CD stages; there is no
actual source checkout, package installation, security scan, or deployment.
An intentional failure is a successful demonstration of the scheduler's behavior.

## Import into the deployed API

Run from the repository root. With curl (on Windows use `curl.exe`):

```sh
curl --fail-with-body -X POST https://flowforge-api-m3ua.onrender.com/api/pipelines -H "Content-Type: application/json" --data-binary @docs/examples/complex-pipeline.json
```

The response is `201` and includes the new pipeline `id`. Open
<https://flowforge-weld-five.vercel.app>, refresh the pipeline list, select the
demo, and click **START RUN**. The event stream shows only events received while
connected. Open a running demo promptly to observe live transitions.

Every import creates another definition. The API currently has no edit/delete
endpoint or import deduplication.

To create and start it with PowerShell:

```powershell
$apiOrigin = 'https://flowforge-api-m3ua.onrender.com'
$body = Get-Content docs/examples/complex-pipeline.json -Raw
$pipeline = Invoke-RestMethod -Method Post -Uri "$apiOrigin/api/pipelines" -ContentType 'application/json; charset=utf-8' -Body ([System.Text.Encoding]::UTF8.GetBytes($body))
$run = Invoke-RestMethod -Method Post -Uri "$apiOrigin/api/pipelines/$($pipeline.id)/runs"
Write-Output "https://flowforge-weld-five.vercel.app/pipelines/$($pipeline.id)"
Write-Output "https://flowforge-weld-five.vercel.app/runs/$($run.id)"
```

## Successful, slower variant

For a full green run with enough time to watch parallel work, transform a copy of
the fixture in memory. This does not change the checked-in failure example.
Run from the repository root in PowerShell:

```powershell
$apiOrigin = 'https://flowforge-api-m3ua.onrender.com'
$demo = Get-Content docs/examples/complex-pipeline.json -Raw | ConvertFrom-Json
$demo.name = 'Platform release — parallel CI/CD demo'
foreach ($job in $demo.jobs) {
    $job.timeout = '30s'
    $job.config.duration = '4s'
}
$unitJob = $demo.jobs | Where-Object id -eq 'unit'
$unitJob.name = 'Unit tests'
$unitJob.config.duration = '7s'
$body = $demo | ConvertTo-Json -Depth 10
$pipeline = Invoke-RestMethod -Method Post -Uri "$apiOrigin/api/pipelines" -ContentType 'application/json; charset=utf-8' -Body ([System.Text.Encoding]::UTF8.GetBytes($body))
Write-Output "https://flowforge-weld-five.vercel.app/pipelines/$($pipeline.id)"
```

Open the printed link and click **START RUN**. This variant should reach
`completed` with all 11 jobs complete. Start another run of this slower variant
and use **CANCEL RUN** while work is active to demonstrate cancellation.

## Retry and replay

On a failed or cancelled run, **RETRY RUN** creates a new run from the same saved
definition. Retrying the intentional-timeout example will reproduce that timeout;
it does not repair the definition. To show success, import the variant above.
To rerun a successful pipeline, return to its definition and use **START RUN**.

Refreshing a run page restores its current snapshot, not past events. The API
does not replay event history or resume old work after a process restart.

## Local API

Start the backend with `go run ./cmd/api` from `backend/`, then import from the
repository root:

```sh
curl --fail-with-body -X POST http://localhost:8080/api/pipelines -H "Content-Type: application/json" --data-binary @docs/examples/complex-pipeline.json
```

For the PowerShell snippets use `$apiOrigin = 'http://localhost:8080'` and open
the frontend at `http://localhost:5173` instead. Local data has the same in-memory
lifetime. See [deployment](../deployment.md) for CORS and environment settings.
