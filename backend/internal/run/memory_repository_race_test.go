package run

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMemoryRepositoryConcurrentSnapshots(t *testing.T) {
	repository := NewMemoryRepository()
	initial := sampleRun("run", time.Now())
	if err := repository.Create(context.Background(), initial); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	const workersCount = 40
	errorsCh := make(chan error, workersCount)
	var workers sync.WaitGroup
	for index := range workersCount {
		workers.Add(1)
		go func() {
			defer workers.Done()
			if index%2 == 0 {
				snapshot := Clone(initial)
				snapshot.Jobs[0].Attempt = index
				if err := repository.Save(context.Background(), snapshot); err != nil {
					errorsCh <- err
				}
				return
			}
			if _, err := repository.GetByID(context.Background(), initial.ID); err != nil {
				errorsCh <- err
			}
		}()
	}
	workers.Wait()
	close(errorsCh)
	for err := range errorsCh {
		t.Error(fmt.Errorf("concurrent repository operation: %w", err))
	}
}
