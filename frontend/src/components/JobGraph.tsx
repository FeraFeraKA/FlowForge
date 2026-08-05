import { Check, Timer, Workflow } from 'lucide-react'
import type { JobDefinition, JobRun } from '@/api/types'
import { StatusBadge } from '@/components/StatusBadge'

type GraphJob = JobDefinition | JobRun
const isRunJob = (job: GraphJob): job is JobRun => 'status' in job

export function JobGraph({ jobs }: { jobs: GraphJob[] }) {
  return (
    <div className="divide-y">
      {jobs.map((job, index) => {
        const running = isRunJob(job)
        return (
          <article key={'id' in job ? job.id : job.jobId} className="relative grid gap-4 p-4 sm:grid-cols-[42px_1fr_auto] sm:items-center">
            <div className="grid size-9 place-items-center border bg-background text-xs text-muted-foreground">
              {running && job.status === 'completed' ? <Check className="size-4 text-emerald-400" /> : String(index + 1).padStart(2, '0')}
            </div>
            <div className="min-w-0">
              <div className="flex items-center gap-3"><h3 className="truncate text-xs font-medium">{job.name}</h3>{running && <StatusBadge status={job.status} pulse={job.status === 'running'} />}</div>
              <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-[9px] text-muted-foreground">
                <span className="flex items-center gap-1"><Workflow className="size-3" />DEP / {job.dependsOn.join(', ') || 'ROOT'}</span>
                <span className="flex items-center gap-1"><Timer className="size-3" />TIMEOUT / {job.timeout ?? 'NONE'}</span>
                <span>RETRIES / {job.maxRetries}</span>
                {running && <span>ATTEMPT / {job.attempt}</span>}
              </div>
              {running && job.status === 'running' && (
                <div className="mt-3 h-1 border bg-muted"><div className="h-full bg-primary transition-all" style={{ width: `${job.progress}%` }} /></div>
              )}
              {running && job.error && <p className="mt-2 text-[9px] text-rose-400">ERR / {job.error}</p>}
            </div>
            <code className="text-[9px] text-muted-foreground">{running ? job.jobId : job.id}</code>
          </article>
        )
      })}
    </div>
  )
}
