# FlowForge frontend

React 19, TypeScript, Vite, React Router, Tailwind CSS, and Radix UI.
The typed REST client and native `EventSource` call the Go API directly.

- App: <https://flowforge-weld-five.vercel.app>
- API: <https://flowforge-api-m3ua.onrender.com>
- [Deployment settings](../docs/deployment.md)
- [Implementation and data flow](../docs/frontend.md)

## Local setup

Use Node.js 24 and pnpm compatible with the version 9 lockfile. From this directory:

```sh
pnpm install --frozen-lockfile
pnpm dev
```

Run `go run ./cmd/api` separately from `backend/`. Open `http://localhost:5173`;
the client defaults to `http://localhost:8080`.

For a different API, create the ignored file `.env.local` in this directory:

```dotenv
VITE_API_URL=https://flowforge-api-m3ua.onrender.com
```

Restart Vite after changes. The API must allow your local browser origin too;
the hosted API currently allows only the primary hosted frontend. Prefer the
local API during development to avoid changing shared service settings.

## Routes

| Route | View |
| --- | --- |
| `/` | Pipeline list |
| `/pipelines/new` | Create a pipeline and jobs |
| `/pipelines/:pipelineId` | Definition and Start Run |
| `/runs/:runId` | Job states, cancel/retry, live events |

Unknown routes redirect to `/`. There is no separate `/pipelines` list route.

## Commands

```sh
pnpm build
pnpm lint
pnpm preview
```

Build type-checks with `tsc -b` and emits `dist/`. Preview uses the API URL
embedded at build time; its origin (usually port 4173) needs an API allowlist entry.

## Vercel

Root: `frontend`; framework: Vite; output: `dist`; Node.js: 24.x.
The configured hosted command is `npm run build`, invoking the same script as
`pnpm build`; installation uses the pnpm lockfile.

Set `VITE_API_URL=https://flowforge-api-m3ua.onrender.com` in Production and
Preview before building. This is public build-time configuration, not a secret.
Render must also allow each frontend origin. `vercel.json` handles direct SPA
links; REST/SSE traffic goes to Render, not a Vercel API route.

## Demo limits

[Import the complex fixture](../docs/examples/README.md) to demonstrate parallel
work, timeouts, retries, and blocked jobs. Demos are not seeded automatically.
Only delay jobs exist; CI/CD stage names are illustrative. The event panel shows
live orchestration events, not shell output or historical logs. Data disappears
when the API process restarts.
