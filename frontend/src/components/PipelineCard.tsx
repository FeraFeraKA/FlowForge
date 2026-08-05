import { ArrowUpRight, Clock3, GitFork, Layers3 } from 'lucide-react'
import { Link } from 'react-router'
import type { Pipeline } from '@/api/types'
import { formatDate } from '@/lib/format'

export function PipelineCard({ pipeline }: { pipeline: Pipeline }) {
  const roots = pipeline.jobs.filter((job) => job.dependsOn.length === 0).length
  return (
    <Link
      to={`/pipelines/${pipeline.id}`}
      className="group relative block border bg-card/80 p-5 transition-colors hover:border-primary/70 hover:bg-accent/30"
    >
      <span className="absolute right-4 top-4 text-muted-foreground transition-colors group-hover:text-primary">
        <ArrowUpRight className="size-4" />
      </span>
      <div className="mb-8 flex size-9 items-center justify-center border bg-background text-primary">
        <GitFork className="size-4" />
      </div>
      <h2 className="pr-8 text-sm font-medium text-foreground">{pipeline.name}</h2>
      <p className="mt-2 truncate text-[10px] text-muted-foreground">ID / {pipeline.id}</p>
      <div className="mt-6 grid grid-cols-3 border-t pt-4 text-[10px] text-muted-foreground">
        <span className="flex items-center gap-1.5"><Layers3 className="size-3" />{pipeline.jobs.length} JOBS</span>
        <span className="flex items-center gap-1.5"><GitFork className="size-3" />{roots} ROOTS</span>
        <span className="flex items-center justify-end gap-1.5"><Clock3 className="size-3" />{formatDate(pipeline.createdAt)}</span>
      </div>
    </Link>
  )
}
