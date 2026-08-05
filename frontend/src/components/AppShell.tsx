import { Activity, Code2, Plus, SquareTerminal, Workflow } from 'lucide-react'
import { useEffect, useState } from 'react'
import { NavLink, Outlet } from 'react-router'
import { api } from '@/api/client'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

const links = [
  { to: '/', label: 'Pipelines', icon: Workflow, end: true },
  { to: '/pipelines/new', label: 'New pipeline', icon: Plus },
]

export function AppShell() {
  const online = useBackendHealth()
  return (
    <div className="min-h-screen bg-background">
      <aside className="fixed inset-y-0 left-0 z-30 hidden w-64 border-r bg-[#0a0a0b] md:flex md:flex-col">
        <NavLink to="/" className="flex h-16 items-center gap-3 border-b px-6 text-sm font-semibold">
          <span className="grid size-7 place-items-center bg-primary text-primary-foreground">
            <SquareTerminal className="size-4" />
          </span>
          FLOWFORGE
        </NavLink>
        <div className="px-3 py-5">
          <p className="mb-3 px-3 text-[10px] tracking-[0.18em] text-muted-foreground">WORKSPACE / LOCAL</p>
          <nav className="space-y-1">
            {links.map(({ to, label, icon: Icon, end }) => (
              <NavLink
                key={to}
                to={to}
                end={end}
                className={({ isActive }) => cn(
                  'flex h-9 items-center gap-3 border border-transparent px-3 text-xs text-muted-foreground transition-colors',
                  isActive && 'border-border bg-accent text-accent-foreground',
                  !isActive && 'hover:border-border hover:text-foreground',
                )}
              >
                <Icon className="size-3.5" />{label}
              </NavLink>
            ))}
          </nav>
        </div>
        <div className="mt-auto border-t p-4 text-[10px] text-muted-foreground">
          <div className={cn('mb-2 flex items-center gap-2', online ? 'text-emerald-400' : 'text-rose-400')}>
            <Activity className="size-3" /> BACKEND {online ? 'CONNECTED' : 'OFFLINE'}
          </div>
          <span>ENV / DEVELOPMENT</span>
        </div>
      </aside>

      <div className="md:pl-64">
        <header className="sticky top-0 z-20 flex h-16 items-center justify-between border-b bg-background/95 px-4 backdrop-blur md:px-8">
          <NavLink to="/" className="flex items-center gap-2 text-xs md:hidden">
            <SquareTerminal className="size-4 text-primary" /> FLOWFORGE
          </NavLink>
          <div className="hidden items-center gap-2 text-[10px] text-muted-foreground md:flex">
            <span className={online ? 'text-emerald-400' : 'text-rose-400'}>●</span> SYSTEM {online ? 'OPERATIONAL' : 'UNAVAILABLE'}
          </div>
          <div className="flex items-center gap-2">
            <Button asChild variant="ghost" size="icon-sm">
              <a href="https://github.com" aria-label="GitHub"><Code2 /></a>
            </Button>
            <Button asChild size="sm"><NavLink to="/pipelines/new"><Plus /> CREATE</NavLink></Button>
          </div>
        </header>
        <main className="dev-grid min-h-[calc(100vh-4rem)]"><Outlet /></main>
      </div>
    </div>
  )
}

function useBackendHealth() {
  const [online, setOnline] = useState(false)
  useEffect(() => {
    let active = true
    const check = () => api.health().then(() => { if (active) setOnline(true) }).catch(() => { if (active) setOnline(false) })
    void check()
    const interval = window.setInterval(check, 10_000)
    return () => { active = false; window.clearInterval(interval) }
  }, [])
  return online
}
