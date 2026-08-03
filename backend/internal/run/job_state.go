package run

func (job *JobRun) TransitionTo(next JobStatus) error {
	if !canJobTransition(job.Status, next) {
		return TransitionError{
			Entity: "job",
			From:   string(job.Status),
			To:     string(next),
		}
	}

	job.Status = next
	return nil
}

func canJobTransition(current, next JobStatus) bool {
	switch current {
	case JobStatusPending:
		return next == JobStatusQueued ||
			next == JobStatusBlocked ||
			next == JobStatusCancelled
	case JobStatusQueued:
		return next == JobStatusRunning ||
			next == JobStatusCancelled
	case JobStatusRunning:
		return next == JobStatusCompleted ||
			next == JobStatusFailed ||
			next == JobStatusCancelled
	case JobStatusFailed:
		return next == JobStatusRetrying
	case JobStatusRetrying:
		return next == JobStatusQueued ||
			next == JobStatusCancelled
	default:
		return false
	}
}
