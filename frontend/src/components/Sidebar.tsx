const navItems = [
  { icon: '⌘', label: 'Overview', active: true },
  { icon: '◇', label: 'Pipelines' },
  { icon: '↻', label: 'Runs' },
  { icon: '⚙', label: 'Settings' },
]

export function Sidebar() {
  return (
    <aside className="sidebar">
      <div className="brand">
        <span className="brand__mark">F</span>
        <span>FlowForge</span>
      </div>
      <nav className="navigation" aria-label="Main navigation">
        {navItems.map((item) => (
          <a className={item.active ? 'nav-item nav-item--active' : 'nav-item'} href="#" key={item.label}>
            <span className="nav-item__icon">{item.icon}</span>
            {item.label}
          </a>
        ))}
      </nav>
      <div className="sidebar__footer">
        <div className="avatar">AK</div>
        <div>
          <strong>Alex Kim</strong>
          <span>Developer</span>
        </div>
        <button className="icon-button" type="button" aria-label="Account options">•••</button>
      </div>
    </aside>
  )
}
