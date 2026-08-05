import { activity } from '../mockData'

export function Activity() {
  return (
    <section className="panel activity-panel">
      <div className="panel__header"><div><h2>Live activity</h2><p>Run #186</p></div><span className="live-label">Live</span></div>
      <div className="activity-list">
        {activity.map((item) => (
          <article className="activity-item" key={item.title}>
            <span className={`activity-item__icon activity-item__icon--${item.tone}`}>
              {item.tone === 'success' ? '✓' : item.tone === 'info' ? '↗' : '▶'}
            </span>
            <div><strong>{item.title}</strong><span>{item.detail}</span></div>
            <time>{item.time}</time>
          </article>
        ))}
      </div>
      <button className="activity-panel__footer" type="button">Open run details <span>→</span></button>
    </section>
  )
}
