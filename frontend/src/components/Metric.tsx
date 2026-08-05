export function Metric({ label, value }: { label: string; value: React.ReactNode }) {
  return <div className="border bg-card/70 p-4"><p className="text-[9px] tracking-wider text-muted-foreground">{label}</p><div className="mt-3 text-sm text-foreground">{value}</div></div>
}
