import { Activity } from './components/Activity'
import { JobList } from './components/JobList'
import { PageHeader } from './components/PageHeader'
import { Sidebar } from './components/Sidebar'
import { Stats } from './components/Stats'
import './App.css'

function App() {
  return (
    <div className="app-shell">
      <Sidebar />
      <main className="main-content">
        <div className="topbar">
          <span className="topbar__workspace"><span>W</span> Web platform <b>⌄</b></span>
          <div className="topbar__actions">
            <button className="icon-button" type="button" aria-label="Search">⌕</button>
            <button className="icon-button notification" type="button" aria-label="Notifications">♢</button>
          </div>
        </div>
        <div className="page">
          <PageHeader />
          <Stats />
          <div className="dashboard-grid">
            <JobList />
            <Activity />
          </div>
        </div>
      </main>
    </div>
  )
}

export default App
