import { GripVertical, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

export type JobDraft = {
  key: number
  id: string
  name: string
  duration: string
  dependsOn: string
  timeout: string
  maxRetries: string
  retryDelay: string
}

export function JobEditor({
  job, index, removable, onChange, onRemove,
}: {
  job: JobDraft
  index: number
  removable: boolean
  onChange: (job: JobDraft) => void
  onRemove: () => void
}) {
  const set = (field: keyof JobDraft, value: string) => onChange({ ...job, [field]: value })
  return (
    <section className="border bg-card/80">
      <header className="flex h-11 items-center justify-between border-b px-3">
        <div className="flex items-center gap-2 text-[10px] tracking-wider text-muted-foreground">
          <GripVertical className="size-3" /> JOB_{String(index + 1).padStart(2, '0')}
        </div>
        <Button type="button" variant="ghost" size="icon-sm" disabled={!removable} onClick={onRemove} aria-label="Remove job">
          <Trash2 className="size-3.5" />
        </Button>
      </header>
      <div className="grid gap-5 p-5 md:grid-cols-2">
        <Field label="JOB ID" hint="tests">
          <Input value={job.id} onChange={(event) => set('id', event.target.value)} placeholder="tests" required />
        </Field>
        <Field label="DISPLAY NAME" hint="Run tests">
          <Input value={job.name} onChange={(event) => set('name', event.target.value)} placeholder="Run tests" required />
        </Field>
        <Field label="DURATION" hint="Go duration, e.g. 2s">
          <Input value={job.duration} onChange={(event) => set('duration', event.target.value)} placeholder="2s" required />
        </Field>
        <Field label="TIMEOUT" hint="Optional, e.g. 10s">
          <Input value={job.timeout} onChange={(event) => set('timeout', event.target.value)} placeholder="10s" />
        </Field>
        <Field label="DEPENDS ON" hint="Comma-separated job IDs">
          <Input value={job.dependsOn} onChange={(event) => set('dependsOn', event.target.value)} placeholder="install, build" />
        </Field>
        <div className="grid grid-cols-2 gap-3">
          <Field label="RETRIES" hint="0 or more">
            <Input type="number" min="0" value={job.maxRetries} onChange={(event) => set('maxRetries', event.target.value)} />
          </Field>
          <Field label="RETRY DELAY" hint="e.g. 500ms">
            <Input value={job.retryDelay} onChange={(event) => set('retryDelay', event.target.value)} />
          </Field>
        </div>
      </div>
    </section>
  )
}

function Field({ label, hint, children }: { label: string; hint: string; children: React.ReactNode }) {
  return <div className="space-y-2"><Label className="text-[10px] tracking-wider">{label}</Label>{children}<p className="text-[9px] text-muted-foreground">{hint}</p></div>
}
