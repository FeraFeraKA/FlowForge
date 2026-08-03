package run

import (
	"errors"
	"testing"
)

func TestRejectedJobTransitions(t *testing.T) {
	statuses := []JobStatus{
		JobStatusPending, JobStatusQueued, JobStatusRunning, JobStatusRetrying,
		JobStatusCompleted, JobStatusFailed, JobStatusBlocked, JobStatusCancelled,
	}
	allowed := map[[2]JobStatus]bool{
		{JobStatusPending, JobStatusQueued}:     true,
		{JobStatusPending, JobStatusBlocked}:    true,
		{JobStatusPending, JobStatusCancelled}:  true,
		{JobStatusQueued, JobStatusRunning}:     true,
		{JobStatusQueued, JobStatusCancelled}:   true,
		{JobStatusRunning, JobStatusCompleted}:  true,
		{JobStatusRunning, JobStatusFailed}:     true,
		{JobStatusRunning, JobStatusCancelled}:  true,
		{JobStatusFailed, JobStatusRetrying}:    true,
		{JobStatusRetrying, JobStatusQueued}:    true,
		{JobStatusRetrying, JobStatusCancelled}: true,
	}

	for _, from := range statuses {
		for _, to := range statuses {
			if allowed[[2]JobStatus{from, to}] {
				continue
			}
			t.Run(string(from)+"_to_"+string(to), func(t *testing.T) {
				job := JobRun{Status: from}
				err := job.TransitionTo(to)
				if !errors.Is(err, ErrInvalidTransition) {
					t.Fatalf("TransitionTo() error = %v, want ErrInvalidTransition", err)
				}
				if job.Status != from {
					t.Errorf("status after rejected transition = %q, want %q", job.Status, from)
				}
			})
		}
	}
}
