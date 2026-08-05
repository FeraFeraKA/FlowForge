package events

import runmodel "flowforge/backend/internal/run"

func changedRunStatus(previous, current runmodel.Run, hadPrevious bool) (Type, bool) {
	if hadPrevious && previous.Status == current.Status {
		return "", false
	}
	switch current.Status {
	case runmodel.RunStatusRunning:
		return RunStarted, true
	case runmodel.RunStatusCompleted:
		return RunCompleted, true
	case runmodel.RunStatusFailed:
		return RunFailed, true
	case runmodel.RunStatusCancelled:
		return RunCancelled, true
	default:
		return "", false
	}
}

func changedJobStatus(previous, current runmodel.JobRun, hadPrevious bool) (Type, bool) {
	if hadPrevious && previous.Status == current.Status {
		return "", false
	}
	switch current.Status {
	case runmodel.JobStatusQueued:
		return JobQueued, true
	case runmodel.JobStatusRunning:
		return JobStarted, true
	case runmodel.JobStatusRetrying:
		return JobRetrying, true
	case runmodel.JobStatusCompleted:
		return JobCompleted, true
	case runmodel.JobStatusFailed:
		return JobFailed, true
	case runmodel.JobStatusBlocked:
		return JobBlocked, true
	case runmodel.JobStatusCancelled:
		return JobCancelled, true
	default:
		return "", false
	}
}
