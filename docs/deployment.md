# Deployment and operations

The Vite frontend is hosted on Vercel and the Go API on Render. The browser calls
Render directly for REST and SSE; Vercel does not proxy the API. These settings
describe the deployment configured on 2026-10-23.

## Live endpoints

- Frontend: <https://flowforge-weld-five.vercel.app>
- API origin: <https://flowforge-api-m3ua.onrender.com>
- Health: <https://flowforge-api-m3ua.onrender.com/health>
- Pipeline list: <https://flowforge-api-m3ua.onrender.com/api/pipelines>
- Source: <https://github.com/FeraFeraKA/FlowForge>, branch `main`.

Use the primary frontend URL above. Other aliases and Vercel preview URLs are
different browser origins and require separate CORS entries.

## Render API

Service: `flowforge-api`, workspace: `My Workspace`.
[Open the service dashboard](https://dashboard.render.com/web/srv-dbdogtm7bikc73eb7cig).

| Setting | Configured value |
| --- | --- |
| Service type | Public web service |
| Runtime | Go |
| Region | Frankfurt |
| Compute plan / instances | Free / 1 |
| Branch | `main` |
| Root directory | Repository root (empty setting) |
| Build command | `cd backend && go test ./... && go build -o bin/api ./cmd/api` |
| Start command | `cd backend && ./bin/api` |
| Auto-deploy | On commit |
| Platform health check | Default TCP; custom HTTP path is unset |

The application independently exposes `GET /health`, returning `200` and
`{"status":"ok"}`. To use it as Render's HTTP probe, set Health Check Path to
`/health` in the dashboard; this is an optional change, not the current setting.

There is no database, persistent disk, migration, pre-deploy command, worker,
or Redis dependency. There is no `render.yaml`; the active service is configured
in Render. The build runs Go tests before producing the binary.

## Vercel frontend

Project: `flowforge`, linked to the same GitHub repository.

| Setting | Value |
| --- | --- |
| Framework | Vite |
| Root directory | `frontend` |
| Node.js | 24.x |
| Dependency lockfile | `frontend/pnpm-lock.yaml` |
| Hosted build command | `npm run build` |
| Build script | `tsc -b && vite build` |
| Output directory | `dist`, relative to `frontend` |

Use `pnpm build` locally. The configured hosted command invokes the same package
script; dependency installation uses the pnpm lockfile. Do not generate an npm
lockfile just to match the script runner.

The project is Git-linked for subsequent deployments from `main`. After a push,
check the resulting deployment's source commit and target in Vercel.
`frontend/vercel.json` provides the SPA fallback for direct application URLs such
as `/pipelines/new` and `/runs/<id>`. It is not an API proxy.

## Environment variables

| Variable | Where | Current value / default | Applied |
| --- | --- | --- | --- |
| `VITE_API_URL` | Vercel Production and Preview | `https://flowforge-api-m3ua.onrender.com` | Frontend build |
| `ALLOWED_ORIGINS` | Render | `https://flowforge-weld-five.vercel.app` | API startup |
| `PORT` | Render-provided | Local fallback: `8080` | API startup |

`VITE_API_URL` is an origin without a trailing slash or `/api` suffix. The client
appends full paths such as `/api/pipelines`. Without it, the client uses
`http://localhost:8080`. Vite embeds this value into public JavaScript; do not put
secrets into it. Changing it requires a new frontend build.

`ALLOWED_ORIGINS` is a comma-separated list of exact origins (scheme, host,
non-default port), without paths or trailing slashes. Whitespace is trimmed.
A non-empty value replaces the local defaults rather than extending them.
Without it, the API permits `http://localhost:5173` and
`http://127.0.0.1:5173`.

Preview URLs, alternate Vercel aliases, and local preview ports must be added
explicitly alongside production, followed by an API deploy with the new environment.
Setting a Preview API URL alone does not authorize that browser origin.
CORS applies to SSE and REST. Requests without an `Origin` header are not
restricted by this middleware: CORS is not authentication.

## Recreate or update hosting

1. Push the intended code to `main`.
2. Configure the Render service with the settings above and the intended frontend
   origin in `ALLOWED_ORIGINS`.
3. Deploy the API, wait for `live`, then check `/health`.
4. Configure Vercel with root `frontend`, framework Vite, and output `dist`.
5. Set `VITE_API_URL` to the actual Render origin before building for Production.
6. Verify the frontend domain assigned by Vercel.
7. If it differs from the Render allowlist, update `ALLOWED_ORIGINS` and deploy
   the API again before testing in the browser.

Changing the API hostname requires rebuilding the frontend; changing the frontend
hostname requires updating the API allowlist. Use provider-assigned domains
rather than deriving them from service names.

This public deployment has no authentication or cookies. Use it for shared demos,
not sensitive workloads. Only delay jobs are supported; no real CI commands run.

## Release checks

Run the [repository checks](../README.md#verification) before releasing code.
After the push:

1. Confirm the intended commit is `READY` in Vercel and `live` in Render.
   A successful deploy trigger alone is not proof of availability.
2. Check the API:

   ```sh
   curl https://flowforge-api-m3ua.onrender.com/health
   curl https://flowforge-api-m3ua.onrender.com/api/pipelines
   ```

3. Open the frontend, create/import a pipeline, and start a run.
4. Confirm job updates and the event stream's `CONNECTED` state.
5. Refresh a pipeline/run URL directly to check SPA routing and snapshot loading.
6. Use the [demo fixture](examples/README.md) to check retry and blocking behavior.
7. Inspect Render build/runtime logs and browser console/network errors.

The initial deployment passed Go package tests and a browser check covering
creation, successful execution, SSE, direct-route reload, and API health.
The connector's Vercel runtime-log query returned `403`, so it did not establish
a clean runtime-log scan. No Vercel log drains were configured at that check;
external monitoring is not configured by this deployment.

## Demo data and restart behavior

The [complex fixture](examples/complex-pipeline.json) is imported manually using
the [demo instructions](examples/README.md); it is not a startup seed.
Each import generates a new pipeline ID. Do not document temporary live demo IDs
as permanent links.

Definitions, runs, and event subscriptions belong to one API process in memory.
Restarting, deploying, replacing, or spinning down that process loses its state.
Active work does not resume; old URLs may return `404`. Re-import the demo and
start another run. Rollback restores code, not lost in-memory data.

Keep one API instance: multiple instances do not share definitions, runs, events,
or cancellation state. Changing the compute plan alone does not add persistence.

## Troubleshooting

| Symptom | Check / action |
| --- | --- |
| Deployed UI calls localhost | Set `VITE_API_URL` for the correct Vercel environment and rebuild. |
| `403 origin_not_allowed` / CORS failure | Compare the exact browser origin with `ALLOWED_ORIGINS`; deploy the corrected environment. |
| Production works but preview fails | Add the preview's exact origin without removing production. |
| Direct page URL has a hosting 404 | Check root `frontend` and the deployed `frontend/vercel.json`. |
| API fails to start | Read Render logs; check the binary path, start command, and `PORT`. |
| First API request is slow | Check the Render service state and allow a stopped free service to start. |
| Old pipeline/run has a 404 | State may have reset; list current definitions and re-import the demo. |
| SSE reconnects / event panel is empty | Check CORS and the stream request; there is no historical replay. Reload for a snapshot. |
| Complex demo fails Unit tests | Expected: three timed-out attempts block Deploy preview. |
| No deployment after a push | Check the linked branch/source commit and provider status; manually deploy if auto-deploy did not start. |

Inspect the first failed build/startup step before retrying. After a rollback or
environment change, recheck API health, CORS, SPA routes, and SSE.
