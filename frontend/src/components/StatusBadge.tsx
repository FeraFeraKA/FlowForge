import type { JobStatus } from '../model'

const labels: Record<JobStatus, string> = {
  completed: 'Completed',
  running: 'Running',
  queued: 'Queued',
  retrying: 'Retrying',
  blocked: 'Blocked',
}

export function StatusBadge({ status }: { status: JobStatus }) {
  return (
    <span className={`status status--${status}`}>
      <span className="status__dot" />
      {labels[status]}
    </span>
  )
}
