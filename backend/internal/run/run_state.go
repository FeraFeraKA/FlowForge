package run

func (r *Run) TransitionTo(next RunStatus) error {
	if !canRunTransition(r.Status, next) {
		return TransitionError{
			Entity: "run",
			From:   string(r.Status),
			To:     string(next),
		}
	}

	r.Status = next
	return nil
}

func canRunTransition(current, next RunStatus) bool {
	switch current {
	case RunStatusPending:
		return next == RunStatusRunning ||
			next == RunStatusCancelled
	case RunStatusRunning:
		return next == RunStatusCompleted ||
			next == RunStatusFailed ||
			next == RunStatusCancelled
	default:
		return false
	}
}
