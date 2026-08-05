export function formatDate(value: string | null | undefined) {
  if (!value) return '—'
  return new Intl.DateTimeFormat('en', {
    month: 'short', day: '2-digit', hour: '2-digit', minute: '2-digit',
  }).format(new Date(value))
}

export function formatDuration(start: string | null, end: string | null) {
  if (!start) return '—'
  const milliseconds = new Date(end ?? Date.now()).getTime() - new Date(start).getTime()
  const seconds = Math.max(0, Math.round(milliseconds / 1000))
  if (seconds < 60) return `${seconds}s`
  return `${Math.floor(seconds / 60)}m ${seconds % 60}s`
}
