package scheduler

import "errors"

func (s *ConcurrentScheduler) stop(runtime *concurrentRuntime, cause error) error {
	runtime.cancel()
	cancelErr := s.base.cancelExecution(runtime.state, cause)
	observeErr := s.notify(runtime.state.execution)
	runtime.workers.Wait()
	return errors.Join(cause, cancelErr, observeErr)
}

func (s *ConcurrentScheduler) complete(runtime *concurrentRuntime) error {
	finishErr := s.base.finishExecution(runtime.state)
	var observeErr error
	if finishErr == nil {
		observeErr = s.notify(runtime.state.execution)
	}
	runtime.cancel()
	runtime.workers.Wait()
	return errors.Join(finishErr, observeErr)
}
