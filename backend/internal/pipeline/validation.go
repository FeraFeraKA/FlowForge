package pipeline

import "strings"

func (p Pipeline) Validate() error {
	errs := make(ValidationErrors, 0)

	if strings.TrimSpace(p.Name) == "" {
		errs = append(errs, ValidationError{
			Field:   "name",
			Code:    "required",
			Message: "pipeline name is required",
		})
	}

	if p.MaxParallel < 1 {
		errs = append(errs, ValidationError{
			Field:   "maxParallel",
			Code:    "must_be_positive",
			Message: "maxParallel must be at least 1",
		})
	}

	if len(p.Jobs) == 0 {
		errs = append(errs, ValidationError{
			Field:   "jobs",
			Code:    "required",
			Message: "pipeline must contain at least one job",
		})
	}

	for index, job := range p.Jobs {
		errs = append(errs, validateJobFields(index, job)...)
	}

	dependencyResult := validateDependencies(p.Jobs)
	errs = append(errs, dependencyResult.errors...)

	if len(errs) == 0 {
		return nil
	}

	return errs
}
