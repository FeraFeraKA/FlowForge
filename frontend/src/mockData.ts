import type { ActivityItem, PipelineJob } from './model'

export const jobs: PipelineJob[] = [
  {
    id: 'install',
    name: 'Install dependencies',
    description: 'Restore and cache project packages',
    status: 'completed',
    duration: '18s',
    dependsOn: [],
    progress: 100,
  },
  {
    id: 'tests',
    name: 'Run test suite',
    description: 'Unit and integration tests',
    status: 'completed',
    duration: '42s',
    dependsOn: ['install'],
    progress: 100,
  },
  {
    id: 'build',
    name: 'Build application',
    description: 'Compile production artifacts',
    status: 'running',
    duration: '1m 06s',
    dependsOn: ['install'],
    progress: 68,
  },
  {
    id: 'deploy',
    name: 'Deploy preview',
    description: 'Publish preview environment',
    status: 'queued',
    duration: '—',
    dependsOn: ['tests', 'build'],
    progress: 0,
  },
]

export const activity: ActivityItem[] = [
  { title: 'Test suite passed', detail: '124 tests completed', time: '12s ago', tone: 'success' },
  { title: 'Build is running', detail: 'Bundling application', time: '38s ago', tone: 'info' },
  { title: 'Run started', detail: 'Triggered manually by you', time: '2m ago', tone: 'warning' },
]
