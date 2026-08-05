import { useCallback, useEffect, useState } from 'react'
import { api } from '@/api/client'
import type { Pipeline } from '@/api/types'

export function usePipelines() {
  const [pipelines, setPipelines] = useState<Pipeline[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  const refresh = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const result = await api.listPipelines()
      setPipelines(result.items)
    } catch (cause) {
      setError(cause instanceof Error ? cause : new Error('Unable to load pipelines'))
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { void refresh() }, [refresh])
  return { pipelines, loading, error, refresh }
}

export function usePipeline(id: string | undefined) {
  const [pipeline, setPipeline] = useState<Pipeline | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    if (!id) return
    let active = true
    api.getPipeline(id)
      .then((result) => { if (active) setPipeline(result) })
      .catch((cause: unknown) => {
        if (active) setError(cause instanceof Error ? cause : new Error('Unable to load pipeline'))
      })
      .finally(() => { if (active) setLoading(false) })
    return () => { active = false }
  }, [id])

  return { pipeline, loading, error }
}
