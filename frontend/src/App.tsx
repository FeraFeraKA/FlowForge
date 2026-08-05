import { Navigate, Route, Routes } from 'react-router'
import { AppShell } from '@/components/AppShell'
import { CreatePipelinePage } from '@/pages/CreatePipelinePage'
import { PipelinePage } from '@/pages/PipelinePage'
import { PipelinesPage } from '@/pages/PipelinesPage'
import { RunPage } from '@/pages/RunPage'

export default function App() {
  return (
    <Routes>
      <Route element={<AppShell />}>
        <Route index element={<PipelinesPage />} />
        <Route path="pipelines/new" element={<CreatePipelinePage />} />
        <Route path="pipelines/:pipelineId" element={<PipelinePage />} />
        <Route path="runs/:runId" element={<RunPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Route>
    </Routes>
  )
}
