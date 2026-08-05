import type {
  ApiErrorBody, CreatePipelineRequest, Pipeline, PipelineList, Run,
} from './types'

export const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

export class ApiError extends Error {
  readonly status: number
  readonly body?: ApiErrorBody

  constructor(status: number, body?: ApiErrorBody) {
    super(body?.error.message ?? `Request failed with status ${status}`)
    this.name = 'ApiError'
    this.status = status
    this.body = body
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, {
    ...init,
    headers: init?.body
      ? { 'Content-Type': 'application/json', ...init.headers }
      : init?.headers,
  })
  if (!response.ok) {
    let body: ApiErrorBody | undefined
    try {
      body = await response.json() as ApiErrorBody
    } catch {
      body = undefined
    }
    throw new ApiError(response.status, body)
  }
  return response.json() as Promise<T>
}

export const api = {
  health: () => request<{ status: string }>('/health'),
  listPipelines: () => request<PipelineList>('/api/pipelines'),
  getPipeline: (id: string) => request<Pipeline>(`/api/pipelines/${id}`),
  createPipeline: (input: CreatePipelineRequest) =>
    request<Pipeline>('/api/pipelines', { method: 'POST', body: JSON.stringify(input) }),
  startRun: (pipelineId: string) =>
    request<Run>(`/api/pipelines/${pipelineId}/runs`, { method: 'POST' }),
  getRun: (runId: string) => request<Run>(`/api/runs/${runId}`),
  cancelRun: (runId: string) =>
    request<Run>(`/api/runs/${runId}/cancel`, { method: 'POST' }),
  retryRun: (runId: string) =>
    request<Run>(`/api/runs/${runId}/retry`, { method: 'POST' }),
}
