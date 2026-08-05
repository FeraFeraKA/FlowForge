package events

import (
	"testing"
	"time"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

func TestObserverPublishesStateChanges(t *testing.T) {
	bus := NewBus(10)
	observer := NewObserver(bus)
	channel, unsubscribe := bus.Subscribe("run")
	defer unsubscribe()
	execution := runmodel.New("run", pipeline.Pipeline{
		ID: "pipeline", Name: "Pipeline", MaxParallel: 1,
		Jobs: []pipeline.JobDefinition{{ID: "job"}},
	}, time.Now())

	execution.Status = runmodel.RunStatusRunning
	if err := observer.Observe(execution); err != nil {
		t.Fatalf("Observe(running) error = %v", err)
	}
	execution.Jobs[0].Status = runmodel.JobStatusQueued
	if err := observer.Observe(execution); err != nil {
		t.Fatalf("Observe(queued) error = %v", err)
	}
	execution.Jobs[0].Status = runmodel.JobStatusRunning
	if err := observer.Observe(execution); err != nil {
		t.Fatalf("Observe(job running) error = %v", err)
	}
	execution.Jobs[0].Status = runmodel.JobStatusCompleted
	execution.Jobs[0].Progress = 100
	if err := observer.Observe(execution); err != nil {
		t.Fatalf("Observe(job completed) error = %v", err)
	}
	execution.Status = runmodel.RunStatusCompleted
	if err := observer.Observe(execution); err != nil {
		t.Fatalf("Observe(completed) error = %v", err)
	}

	want := []Type{RunStarted, JobQueued, JobStarted, JobCompleted, RunCompleted}
	for _, wantType := range want {
		if event := <-channel; event.Type != wantType {
			t.Errorf("event type = %q, want %q", event.Type, wantType)
		}
	}
}
