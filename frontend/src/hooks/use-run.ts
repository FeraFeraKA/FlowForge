import { useCallback, useEffect, useState } from 'react'
import { API_URL, api } from '@/api/client'
import type { Run, RunEvent } from '@/api/types'

const eventTypes = [
  'run.started', 'run.completed', 'run.failed', 'run.cancelled',
  'job.queued', 'job.started', 'job.progress', 'job.retrying',
  'job.completed', 'job.failed', 'job.blocked', 'job.cancelled',
]

function applyEvent(current: Run | null, event: RunEvent): Run | null {
  if (event.payload.run) return event.payload.run
  if (!current || !event.payload.job) return current
  const changed = event.payload.job
  return {
    ...current,
    jobs: current.jobs.map((job) =>
      job.jobId === changed.jobId ? { ...job, ...changed } : job,
    ),
  }
}

export function useRun(runId: string | null) {
  const [run, setRun] = useState<Run | null>(null)
  const [loading, setLoading] = useState(Boolean(runId))
  const [error, setError] = useState<Error | null>(null)
  const [connected, setConnected] = useState(false)

  const refresh = useCallback(async () => {
    if (!runId) return
    try {
      setRun(await api.getRun(runId))
      setError(null)
    } catch (cause) {
      setError(cause instanceof Error ? cause : new Error('Unable to load run'))
    } finally {
      setLoading(false)
    }
  }, [runId])

  useEffect(() => {
    if (!runId) { setRun(null); setLoading(false); return }
    let active = true
    let source: EventSource | null = null
    let retry: number | undefined

    const connect = async () => {
      await refresh()
      if (!active) return
      source = new EventSource(`${API_URL}/api/runs/${runId}/events`)
      source.onopen = () => { if (active) setConnected(true) }
      const receive = (message: MessageEvent<string>) => {
        try {
          const event = JSON.parse(message.data) as RunEvent
          setRun((current) => applyEvent(current, event))
        } catch { setError(new Error('Received an invalid run event')) }
      }
      eventTypes.forEach((type) => source?.addEventListener(type, receive as EventListener))
      source.onerror = () => {
        source?.close()
        setConnected(false)
        if (active) retry = window.setTimeout(() => { void connect() }, 2000)
      }
    }

    void connect()
    return () => {
      active = false
      source?.close()
      if (retry) window.clearTimeout(retry)
    }
  }, [refresh, runId])

  return { run, loading, error, connected, refresh }
}
