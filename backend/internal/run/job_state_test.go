package run

import (
	"errors"
	"testing"
)

func TestAllowedJobTransitions(t *testing.T) {
	transitions := [][2]JobStatus{
		{JobStatusPending, JobStatusQueued},
		{JobStatusPending, JobStatusBlocked},
		{JobStatusPending, JobStatusCancelled},
		{JobStatusQueued, JobStatusRunning},
		{JobStatusQueued, JobStatusCancelled},
		{JobStatusRunning, JobStatusCompleted},
		{JobStatusRunning, JobStatusFailed},
		{JobStatusRunning, JobStatusCancelled},
		{JobStatusFailed, JobStatusRetrying},
		{JobStatusRetrying, JobStatusQueued},
		{JobStatusRetrying, JobStatusCancelled},
	}

	for _, transition := range transitions {
		from, to := transition[0], transition[1]
		t.Run(string(from)+"_to_"+string(to), func(t *testing.T) {
			job := JobRun{Status: from}
			if err := job.TransitionTo(to); err != nil {
				t.Fatalf("TransitionTo() error = %v", err)
			}
			if job.Status != to {
				t.Errorf("status = %q, want %q", job.Status, to)
			}
		})
	}
}

func TestJobTransitionErrorContainsContext(t *testing.T) {
	job := JobRun{Status: JobStatusPending}
	err := job.TransitionTo(JobStatusCompleted)

	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("TransitionTo() error = %v, want ErrInvalidTransition", err)
	}
	var transitionErr TransitionError
	if !errors.As(err, &transitionErr) {
		t.Fatalf("TransitionTo() error type = %T, want TransitionError", err)
	}
	if transitionErr.Entity != "job" || transitionErr.From != "pending" || transitionErr.To != "completed" {
		t.Errorf("TransitionError = %#v, want job pending to completed", transitionErr)
	}
	if job.Status != JobStatusPending {
		t.Errorf("status after rejected transition = %q, want pending", job.Status)
	}
}
