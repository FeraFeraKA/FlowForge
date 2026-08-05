export type JobStatus =
  | 'completed'
  | 'running'
  | 'queued'
  | 'retrying'
  | 'blocked'

export type PipelineJob = {
  id: string
  name: string
  description: string
  status: JobStatus
  duration: string
  dependsOn: string[]
  progress: number
}

export type ActivityItem = {
  title: string
  detail: string
  time: string
  tone: 'success' | 'info' | 'warning'
}
