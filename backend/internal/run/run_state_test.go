package run

import (
	"errors"
	"testing"
)

func TestRunTransitions(t *testing.T) {
	statuses := []RunStatus{
		RunStatusPending, RunStatusRunning, RunStatusCompleted,
		RunStatusFailed, RunStatusCancelled,
	}
	allowed := map[[2]RunStatus]bool{
		{RunStatusPending, RunStatusRunning}:   true,
		{RunStatusPending, RunStatusCancelled}: true,
		{RunStatusRunning, RunStatusCompleted}: true,
		{RunStatusRunning, RunStatusFailed}:    true,
		{RunStatusRunning, RunStatusCancelled}: true,
	}

	for _, from := range statuses {
		for _, to := range statuses {
			t.Run(string(from)+"_to_"+string(to), func(t *testing.T) {
				run := Run{Status: from}
				err := run.TransitionTo(to)
				if allowed[[2]RunStatus{from, to}] {
					if err != nil || run.Status != to {
						t.Fatalf("TransitionTo() = (%q, %v), want (%q, nil)", run.Status, err, to)
					}
					return
				}

				if !errors.Is(err, ErrInvalidTransition) {
					t.Fatalf("TransitionTo() error = %v, want ErrInvalidTransition", err)
				}
				if run.Status != from {
					t.Errorf("status after rejected transition = %q, want %q", run.Status, from)
				}
			})
		}
	}
}

func TestRunStatusIsTerminal(t *testing.T) {
	tests := map[RunStatus]bool{
		RunStatusPending: false, RunStatusRunning: false,
		RunStatusCompleted: true, RunStatusFailed: true, RunStatusCancelled: true,
	}
	for status, want := range tests {
		if got := status.IsTerminal(); got != want {
			t.Errorf("%q.IsTerminal() = %t, want %t", status, got, want)
		}
	}
}
