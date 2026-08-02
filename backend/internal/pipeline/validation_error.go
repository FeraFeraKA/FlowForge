package pipeline

import "fmt"

type ValidationError struct {
	Field   string
	Code    string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}

	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	switch len(e) {
	case 0:
		return "pipeline validation failed"
	case 1:
		return e[0].Error()
	default:
		return fmt.Sprintf("%s (and %d more validation errors)", e[0], len(e)-1)
	}
}
