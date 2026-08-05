import { AlertTriangle, Plus, RefreshCw } from 'lucide-react'
import { Link } from 'react-router'
import { PipelineCard } from '@/components/PipelineCard'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { usePipelines } from '@/hooks/use-pipelines'

export function PipelinesPage() {
  const { pipelines, loading, error, refresh } = usePipelines()
  return (
    <div className="mx-auto max-w-7xl px-4 py-10 md:px-8 md:py-14">
      <div className="mb-10 flex flex-col gap-5 border-b pb-8 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="mb-3 text-[10px] tracking-[0.2em] text-primary">ORCHESTRATION / PIPELINES</p>
          <h1 className="text-2xl font-medium tracking-tight md:text-3xl">Execution pipelines</h1>
          <p className="mt-3 max-w-xl text-xs leading-5 text-muted-foreground">Define dependency graphs, control concurrency and observe every job in real time.</p>
        </div>
        <Button asChild><Link to="/pipelines/new"><Plus /> NEW PIPELINE</Link></Button>
      </div>

      {loading && <div className="grid gap-4 lg:grid-cols-2">{[1, 2, 3, 4].map((item) => <Skeleton key={item} className="h-52 border" />)}</div>}
      {error && (
        <div className="border border-rose-500/50 bg-rose-500/5 p-6 text-xs text-rose-300">
          <AlertTriangle className="mb-3 size-5" />{error.message}
          <Button variant="outline" size="sm" className="mt-5 block" onClick={() => void refresh()}><RefreshCw /> RETRY</Button>
        </div>
      )}
      {!loading && !error && pipelines.length === 0 && (
        <div className="border bg-card/70 px-6 py-20 text-center">
          <pre className="terminal-glow mb-6 text-xs leading-5 text-primary">{`┌──────────────────────┐\n│  NO PIPELINES FOUND  │\n└──────────────────────┘`}</pre>
          <p className="mb-6 text-xs text-muted-foreground">Create the first graph to start orchestrating work.</p>
          <Button asChild><Link to="/pipelines/new"><Plus /> CREATE PIPELINE</Link></Button>
        </div>
      )}
      {!loading && !error && pipelines.length > 0 && (
        <div className="grid gap-4 lg:grid-cols-2">{pipelines.map((pipeline) => <PipelineCard key={pipeline.id} pipeline={pipeline} />)}</div>
      )}
    </div>
  )
}
