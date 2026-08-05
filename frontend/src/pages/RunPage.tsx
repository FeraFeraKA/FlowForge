import { ArrowLeft, Ban, Clock3, RotateCcw, Workflow } from 'lucide-react'
import { Link, useNavigate, useParams } from 'react-router'
import { useState } from 'react'
import { api } from '@/api/client'
import { EventLog } from '@/components/EventLog'
import { JobGraph } from '@/components/JobGraph'
import { Metric } from '@/components/Metric'
import { StatusBadge } from '@/components/StatusBadge'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { useRun } from '@/hooks/use-run'
import { formatDate, formatDuration } from '@/lib/format'

const terminal = new Set(['completed', 'failed', 'cancelled'])

export function RunPage() {
  const { runId } = useParams()
  const navigate = useNavigate()
  const { run, events, loading, error, connected } = useRun(runId ?? null)
  const [command, setCommand] = useState<'cancel' | 'retry' | null>(null)
  const [commandError, setCommandError] = useState<Error | null>(null)

  const cancel = async () => {
    if (!run) return
    setCommand('cancel'); setCommandError(null)
    try { await api.cancelRun(run.id) }
    catch (cause) { setCommandError(cause instanceof Error ? cause : new Error('Unable to cancel run')) }
    finally { setCommand(null) }
  }

  const retry = async () => {
    if (!run) return
    setCommand('retry'); setCommandError(null)
    try {
      const next = await api.retryRun(run.id)
      void navigate(`/runs/${next.id}`)
    } catch (cause) {
      setCommandError(cause instanceof Error ? cause : new Error('Unable to retry run'))
      setCommand(null)
    }
  }

  if (loading) return <div className="mx-auto max-w-6xl space-y-4 px-4 py-14"><Skeleton className="h-28 border" /><Skeleton className="h-96 border" /></div>
  if (error || !run) return <div className="grid min-h-[70vh] place-items-center text-center"><div><p className="text-sm text-rose-400">RUN UNAVAILABLE</p><p className="mt-3 text-xs text-muted-foreground">{error?.message}</p><Button asChild variant="outline" className="mt-6"><Link to="/">PIPELINES</Link></Button></div></div>

  const isTerminal = terminal.has(run.status)
  return (
    <div className="mx-auto max-w-6xl px-4 py-10 md:px-8 md:py-14">
      <Link to={`/pipelines/${run.pipelineId}`} className="mb-8 inline-flex items-center gap-2 text-[10px] text-muted-foreground hover:text-primary"><ArrowLeft className="size-3" /> {run.pipelineName}</Link>
      <header className="flex flex-col gap-6 border-b pb-8 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <div className="mb-3 flex items-center gap-3"><p className="text-[10px] tracking-[0.18em] text-primary">RUN / {run.id}</p><StatusBadge status={run.status} pulse={run.status === 'running'} /></div>
          <h1 className="text-2xl font-medium md:text-3xl">{run.pipelineName}</h1>
          <p className="mt-3 text-xs text-muted-foreground">CREATED / {formatDate(run.createdAt)}</p>
        </div>
        <div className="flex gap-2">
          {!isTerminal && <Button variant="destructive" onClick={() => void cancel()} disabled={command !== null}><Ban />{command === 'cancel' ? 'CANCELLING...' : 'CANCEL RUN'}</Button>}
          {(run.status === 'failed' || run.status === 'cancelled') && <Button onClick={() => void retry()} disabled={command !== null}><RotateCcw />{command === 'retry' ? 'RETRYING...' : 'RETRY RUN'}</Button>}
        </div>
      </header>

      {commandError && <div className="mt-5 border border-rose-500/50 bg-rose-500/5 p-4 text-[10px] text-rose-300">ERR / {commandError.message}</div>}
      <section className="my-6 grid grid-cols-2 gap-3 md:grid-cols-4">
        <Metric label="STATUS" value={<StatusBadge status={run.status} />} />
        <Metric label="DURATION" value={<span className="flex items-center gap-2"><Clock3 className="size-3 text-primary" />{formatDuration(run.startedAt, run.finishedAt)}</span>} />
        <Metric label="JOBS COMPLETE" value={`${run.jobs.filter((job) => job.status === 'completed').length} / ${run.jobs.length}`} />
        <Metric label="MAX PARALLEL" value={run.maxParallel} />
      </section>

      <div className="grid gap-5 lg:grid-cols-[1.35fr_.65fr]">
        <section className="border bg-card/80"><header className="flex h-12 items-center justify-between border-b px-4"><span className="flex items-center gap-2 text-[10px]"><Workflow className="size-3 text-primary" /> JOB EXECUTION</span><span className="text-[9px] text-muted-foreground">COORDINATOR / ACTIVE</span></header><JobGraph jobs={run.jobs} /></section>
        <EventLog events={events} connected={connected} />
      </div>
    </div>
  )
}
