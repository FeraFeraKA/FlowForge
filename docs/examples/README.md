# Demo pipelines

`complex-pipeline.json` demonstrates two root jobs, a dependency fan-out,
`MaxParallel`, successful independent branches, a job timeout, two retries, and
a downstream job becoming blocked after the final failed attempt.

Create it against a local backend:

```bash
curl -X POST http://localhost:8080/api/pipelines \
  -H 'Content-Type: application/json' \
  --data @docs/examples/complex-pipeline.json
```

Start it from the FlowForge UI. To demonstrate cancellation separately, start a
second Run and cancel it while jobs are active; the original Run remains intact.
