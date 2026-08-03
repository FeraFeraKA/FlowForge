package run

import (
	"errors"
	"fmt"
)

var ErrInvalidTransition = errors.New("invalid state transition")

type TransitionError struct {
	Entity string
	From   string
	To     string
}

func (e TransitionError) Error() string {
	return fmt.Sprintf(
		"%s cannot transition from %q to %q",
		e.Entity,
		e.From,
		e.To,
	)
}

func (e TransitionError) Unwrap() error {
	return ErrInvalidTransition
}
