import { Radio } from 'lucide-react'
import type { RunEvent } from '@/api/types'
import { ScrollArea } from '@/components/ui/scroll-area'

export function EventLog({ events, connected }: { events: RunEvent[]; connected: boolean }) {
  return (
    <section className="border bg-[#09090a]">
      <header className="flex h-12 items-center justify-between border-b px-4">
        <span className="flex items-center gap-2 text-[10px]"><Radio className="size-3 text-primary" /> EVENT STREAM</span>
        <span className={connected ? 'text-[9px] text-emerald-400' : 'text-[9px] text-amber-400'}>{connected ? '● CONNECTED' : '○ RECONNECTING'}</span>
      </header>
      <ScrollArea className="h-56">
        <div className="space-y-1 p-4 font-mono text-[9px] leading-5">
          {events.length === 0 && <p className="text-muted-foreground">$ waiting for realtime events<span className="animate-pulse">_</span></p>}
          {events.map((event) => (
            <p key={`${event.id}-${event.type}`} className="grid grid-cols-[72px_1fr] gap-3">
              <time className="text-zinc-600">{new Date(event.timestamp).toLocaleTimeString()}</time>
              <span><b className="font-normal text-primary">{event.type}</b>{event.jobId && <span className="text-muted-foreground"> / {event.jobId}</span>}</span>
            </p>
          ))}
        </div>
      </ScrollArea>
    </section>
  )
}
