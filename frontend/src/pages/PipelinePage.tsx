import { ArrowLeft, GitFork, Play, RotateCw } from 'lucide-react'
import { Link, useNavigate, useParams } from 'react-router'
import { useState } from 'react'
import { api } from '@/api/client'
import { JobGraph } from '@/components/JobGraph'
import { Metric } from '@/components/Metric'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { usePipeline } from '@/hooks/use-pipelines'
import { formatDate } from '@/lib/format'

export function PipelinePage() {
  const { pipelineId } = useParams()
  const navigate = useNavigate()
  const { pipeline, loading, error } = usePipeline(pipelineId)
  const [starting, setStarting] = useState(false)
  const [commandError, setCommandError] = useState<Error | null>(null)

  const start = async () => {
    if (!pipeline) return
    setStarting(true)
    setCommandError(null)
    try {
      const run = await api.startRun(pipeline.id)
      void navigate(`/runs/${run.id}`)
    } catch (cause) {
      setCommandError(cause instanceof Error ? cause : new Error('Unable to start run'))
      setStarting(false)
    }
  }

  if (loading) return <div className="mx-auto max-w-6xl space-y-4 px-4 py-14"><Skeleton className="h-24 border" /><Skeleton className="h-96 border" /></div>
  if (error || !pipeline) return <Message title="PIPELINE UNAVAILABLE" detail={error?.message ?? 'Pipeline not found'} />

  return (
    <div className="mx-auto max-w-6xl px-4 py-10 md:px-8 md:py-14">
      <Link to="/" className="mb-8 inline-flex items-center gap-2 text-[10px] text-muted-foreground hover:text-primary"><ArrowLeft className="size-3" /> PIPELINES</Link>
      <header className="flex flex-col gap-6 border-b pb-8 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="mb-3 text-[10px] tracking-[0.2em] text-primary">PIPELINE / {pipeline.id}</p>
          <h1 className="text-2xl font-medium md:text-3xl">{pipeline.name}</h1>
          <p className="mt-3 text-xs text-muted-foreground">READY FOR EXECUTION / DELAY EXECUTOR</p>
        </div>
        <Button onClick={() => void start()} disabled={starting}><Play />{starting ? 'ALLOCATING...' : 'START RUN'}</Button>
      </header>

      {commandError && <div className="mt-5 border border-rose-500/50 bg-rose-500/5 p-4 text-[10px] text-rose-300">ERR / {commandError.message}</div>}
      <section className="my-6 grid grid-cols-2 gap-3 md:grid-cols-4">
        <Metric label="JOBS" value={pipeline.jobs.length} />
        <Metric label="MAX PARALLEL" value={pipeline.maxParallel} />
        <Metric label="ROOT NODES" value={pipeline.jobs.filter((job) => job.dependsOn.length === 0).length} />
        <Metric label="CREATED" value={formatDate(pipeline.createdAt)} />
      </section>

      <section className="border bg-card/80">
        <div className="flex items-center justify-between border-b p-4">
          <div className="flex items-center gap-2 text-xs"><GitFork className="size-4 text-primary" /> DEPENDENCY GRAPH</div>
          <span className="text-[9px] text-muted-foreground">DAG / VALIDATED</span>
        </div>
        <JobGraph jobs={pipeline.jobs} />
      </section>
    </div>
  )
}

function Message({ title, detail }: { title: string; detail: string }) {
  return <div className="grid min-h-[70vh] place-items-center px-4 text-center"><div><RotateCw className="mx-auto mb-4 size-5 text-primary" /><h1 className="text-sm">{title}</h1><p className="mt-2 text-xs text-muted-foreground">{detail}</p><Button asChild variant="outline" className="mt-6"><Link to="/">BACK TO PIPELINES</Link></Button></div></div>
}
