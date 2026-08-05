import { Badge } from '@/components/ui/badge'
import { cn } from '@/lib/utils'

const tones: Record<string, string> = {
  completed: 'border-emerald-500/50 bg-emerald-500/10 text-emerald-400',
  running: 'border-primary/60 bg-primary/10 text-primary',
  queued: 'border-zinc-600 bg-zinc-800/50 text-zinc-400',
  pending: 'border-zinc-600 bg-zinc-800/50 text-zinc-400',
  retrying: 'border-amber-500/50 bg-amber-500/10 text-amber-400',
  failed: 'border-rose-500/50 bg-rose-500/10 text-rose-400',
  blocked: 'border-orange-500/50 bg-orange-500/10 text-orange-400',
  cancelled: 'border-zinc-500 bg-zinc-800 text-zinc-400',
}

export function StatusBadge({ status, pulse = false }: { status: string; pulse?: boolean }) {
  return (
    <Badge variant="outline" className={cn('gap-1.5 px-2 py-0.5 text-[9px] uppercase tracking-wider', tones[status])}>
      <span className={cn('size-1.5 bg-current', pulse && 'animate-pulse')} />
      {status}
    </Badge>
  )
}
