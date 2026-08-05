import { jobs } from '../mockData'
import type { PipelineJob } from '../model'
import { StatusBadge } from './StatusBadge'

function JobRow({ job, index }: { job: PipelineJob; index: number }) {
  return (
    <article className="job-row">
      <div className={`job-node job-node--${job.status}`}>
        {job.status === 'completed' ? '✓' : index + 1}
      </div>
      <div className="job-row__body">
        <div className="job-row__title">
          <div><strong>{job.name}</strong><span>{job.description}</span></div>
          <StatusBadge status={job.status} />
        </div>
        <div className="job-row__meta">
          <span><b>ID</b> {job.id}</span>
          <span><b>Depends on</b> {job.dependsOn.join(', ') || 'None'}</span>
          <span className="job-row__duration">◷ {job.duration}</span>
        </div>
        {job.status === 'running' && (
          <div className="progress" aria-label={`${job.progress}% complete`}>
            <span style={{ width: `${job.progress}%` }} />
          </div>
        )}
      </div>
    </article>
  )
}

export function JobList() {
  return (
    <section className="panel pipeline-panel">
      <div className="panel__header">
        <div><h2>Pipeline jobs</h2><p>4 jobs · Maximum 2 in parallel</p></div>
        <button className="text-button" type="button">View configuration</button>
      </div>
      <div className="job-list">
        {jobs.map((job, index) => <JobRow job={job} index={index} key={job.id} />)}
      </div>
    </section>
  )
}
