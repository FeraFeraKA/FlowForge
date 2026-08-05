const stats = [
  { label: 'Success rate', value: '94.2%', detail: '+2.4% this week', tone: 'positive' },
  { label: 'Average duration', value: '2m 18s', detail: '12s faster', tone: 'positive' },
  { label: 'Total runs', value: '186', detail: '23 this week', tone: 'neutral' },
]

export function Stats() {
  return (
    <section className="stats" aria-label="Pipeline statistics">
      {stats.map((stat) => (
        <article className="stat-card" key={stat.label}>
          <span className="stat-card__label">{stat.label}</span>
          <strong>{stat.value}</strong>
          <span className={`stat-card__detail stat-card__detail--${stat.tone}`}>{stat.detail}</span>
        </article>
      ))}
      <article className="stat-card stat-card--status">
        <span className="stat-card__label">Current status</span>
        <strong><span className="live-dot" /> Running</strong>
        <span className="stat-card__detail">Run #186 · 1m 24s</span>
      </article>
    </section>
  )
}
