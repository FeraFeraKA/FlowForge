export type RunStatus = 'pending' | 'running' | 'completed' | 'failed' | 'cancelled'
export type JobStatus =
  | 'pending' | 'queued' | 'running' | 'retrying'
  | 'completed' | 'failed' | 'blocked' | 'cancelled'

export type DelayConfig = { duration: string }

export type JobDefinition = {
  id: string
  name: string
  type: 'delay'
  dependsOn: string[]
  timeout?: string
  maxRetries: number
  retryDelay: string
  config?: DelayConfig
}

export type Pipeline = {
  id: string
  name: string
  maxParallel: number
  jobs: JobDefinition[]
  createdAt: string
}

export type PipelineList = { items: Pipeline[]; total: number }
export type CreatePipelineRequest = Omit<Pipeline, 'id' | 'createdAt'>

export type JobRun = {
  jobId: string
  name: string
  type: string
  dependsOn: string[]
  timeout?: string
  maxRetries: number
  retryDelay: string
  config?: DelayConfig
  status: JobStatus
  attempt: number
  progress: number
  error: string | null
  startedAt: string | null
  finishedAt: string | null
}

export type Run = {
  id: string
  pipelineId: string
  pipelineName: string
  status: RunStatus
  maxParallel: number
  jobs: JobRun[]
  createdAt: string
  startedAt: string | null
  finishedAt: string | null
}

export type ErrorDetail = { field: string; code: string; message: string }
export type ApiErrorBody = {
  error: { code: string; message: string; details?: ErrorDetail[] }
}

export type RunEvent = {
  id: string
  type: string
  runId: string
  jobId?: string
  timestamp: string
  payload: { run?: Run; job?: Partial<JobRun> & Pick<JobRun, 'jobId' | 'status'> }
}
