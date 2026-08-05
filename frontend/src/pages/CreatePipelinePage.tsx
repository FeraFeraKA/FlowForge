import { useState } from 'react'
import { ArrowLeft, Plus, Save } from 'lucide-react'
import { Link, useNavigate } from 'react-router'
import { ApiError, api } from '@/api/client'
import type { CreatePipelineRequest } from '@/api/types'
import { JobEditor, type JobDraft } from '@/components/JobEditor'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

let nextKey = 1
const newJob = (): JobDraft => ({
  key: nextKey++, id: '', name: '', duration: '2s', dependsOn: '',
  timeout: '10s', maxRetries: '0', retryDelay: '500ms',
})

function toRequest(name: string, maxParallel: string, jobs: JobDraft[]): CreatePipelineRequest {
  return {
    name,
    maxParallel: Number(maxParallel),
    jobs: jobs.map((job) => ({
      id: job.id.trim(), name: job.name.trim(), type: 'delay',
      dependsOn: job.dependsOn.split(',').map((id) => id.trim()).filter(Boolean),
      timeout: job.timeout.trim() || undefined,
      maxRetries: Number(job.maxRetries), retryDelay: job.retryDelay.trim() || '0s',
      config: { duration: job.duration.trim() },
    })),
  }
}

export function CreatePipelinePage() {
  const navigate = useNavigate()
  const [name, setName] = useState('')
  const [maxParallel, setMaxParallel] = useState('2')
  const [jobs, setJobs] = useState<JobDraft[]>([newJob()])
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<ApiError | Error | null>(null)

  const submit = async (event: React.FormEvent) => {
    event.preventDefault()
    setSubmitting(true)
    setError(null)
    try {
      const pipeline = await api.createPipeline(toRequest(name, maxParallel, jobs))
      void navigate(`/pipelines/${pipeline.id}`)
    } catch (cause) {
      setError(cause instanceof Error ? cause : new Error('Unable to create pipeline'))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <form onSubmit={(event) => void submit(event)} className="mx-auto max-w-5xl px-4 py-10 md:px-8 md:py-14">
      <Link to="/" className="mb-8 inline-flex items-center gap-2 text-[10px] text-muted-foreground hover:text-primary"><ArrowLeft className="size-3" /> PIPELINES</Link>
      <div className="mb-9 border-b pb-7">
        <p className="mb-3 text-[10px] tracking-[0.2em] text-primary">CONFIG / NEW PIPELINE</p>
        <h1 className="text-2xl font-medium">Define execution graph</h1>
        <p className="mt-3 text-xs text-muted-foreground">Jobs use the safe Delay executor. Dependencies are validated by the backend.</p>
      </div>

      <section className="mb-5 grid gap-5 border bg-card/80 p-5 md:grid-cols-[1fr_180px]">
        <div className="space-y-2"><Label>PIPELINE NAME</Label><Input value={name} onChange={(event) => setName(event.target.value)} placeholder="Build and deploy" required /></div>
        <div className="space-y-2"><Label>MAX PARALLEL</Label><Input type="number" min="1" value={maxParallel} onChange={(event) => setMaxParallel(event.target.value)} required /></div>
      </section>

      <div className="space-y-4">
        {jobs.map((job, index) => (
          <JobEditor
            key={job.key} job={job} index={index} removable={jobs.length > 1}
            onChange={(changed) => setJobs((current) => current.map((item) => item.key === job.key ? changed : item))}
            onRemove={() => setJobs((current) => current.filter((item) => item.key !== job.key))}
          />
        ))}
      </div>
      <Button type="button" variant="outline" className="mt-4 w-full border-dashed" onClick={() => setJobs((current) => [...current, newJob()])}><Plus /> ADD JOB</Button>

      {error && <ErrorPanel error={error} />}
      <div className="mt-7 flex justify-end gap-3 border-t pt-6">
        <Button asChild variant="ghost"><Link to="/">CANCEL</Link></Button>
        <Button type="submit" disabled={submitting}><Save />{submitting ? 'SAVING...' : 'SAVE PIPELINE'}</Button>
      </div>
    </form>
  )
}

function ErrorPanel({ error }: { error: ApiError | Error }) {
  const details = error instanceof ApiError ? error.body?.error.details : undefined
  return <div className="mt-5 border border-rose-500/50 bg-rose-500/5 p-4 text-[10px] text-rose-300"><p>{error.message}</p>{details?.map((detail) => <p className="mt-2" key={`${detail.field}-${detail.code}`}>{detail.field} / {detail.message}</p>)}</div>
}
