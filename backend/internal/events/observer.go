package events

import (
	"sync"

	runmodel "flowforge/backend/internal/run"
)

type RunPayload struct {
	Run runmodel.Run `json:"run"`
}

type JobPayload struct {
	Job runmodel.JobRun `json:"job"`
}

type Observer struct {
	mu       sync.Mutex
	bus      *Bus
	previous map[string]runmodel.Run
}

func NewObserver(bus *Bus) *Observer {
	return &Observer{bus: bus, previous: make(map[string]runmodel.Run)}
}

func (o *Observer) Observe(snapshot runmodel.Run) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	previous, exists := o.previous[snapshot.ID]
	if eventType, changed := changedRunStatus(previous, snapshot, exists); changed {
		o.bus.Publish(Event{
			Type: eventType, RunID: snapshot.ID,
			Payload: RunPayload{Run: runmodel.Clone(snapshot)},
		})
	}
	previousJobs := make(map[string]runmodel.JobRun, len(previous.Jobs))
	for _, job := range previous.Jobs {
		previousJobs[job.JobID] = job
	}
	for _, job := range snapshot.Jobs {
		old, hadOld := previousJobs[job.JobID]
		if eventType, changed := changedJobStatus(old, job, hadOld); changed {
			o.bus.Publish(Event{
				Type: eventType, RunID: snapshot.ID, JobID: job.JobID,
				Payload: JobPayload{Job: job},
			})
			continue
		}
		if hadOld && old.Progress != job.Progress {
			o.bus.Publish(Event{
				Type: JobProgress, RunID: snapshot.ID, JobID: job.JobID,
				Payload: JobPayload{Job: job},
			})
		}
	}
	o.previous[snapshot.ID] = runmodel.Clone(snapshot)
	return nil
}
