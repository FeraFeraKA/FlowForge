package pipeline

import "fmt"

type dependencyValidationResult struct {
	errors     ValidationErrors
	graphValid bool
}

func validateDependencies(jobs []JobDefinition) dependencyValidationResult {
	result := dependencyValidationResult{
		errors:     make(ValidationErrors, 0),
		graphValid: true,
	}
	jobIDs := make(map[string]struct{}, len(jobs))

	for index, job := range jobs {
		if job.ID == "" || !validJobID(job.ID) {
			result.graphValid = false
			continue
		}
		if _, exists := jobIDs[job.ID]; exists {
			result.errors = append(result.errors, ValidationError{
				Field:   fmt.Sprintf("jobs[%d].id", index),
				Code:    "duplicate",
				Message: fmt.Sprintf("job ID %q is used more than once", job.ID),
			})
			result.graphValid = false
			continue
		}
		jobIDs[job.ID] = struct{}{}
	}

	for jobIndex, job := range jobs {
		seenDependencies := make(map[string]struct{}, len(job.DependsOn))
		for dependencyIndex, dependencyID := range job.DependsOn {
			field := fmt.Sprintf("jobs[%d].dependsOn[%d]", jobIndex, dependencyIndex)
			if _, duplicate := seenDependencies[dependencyID]; duplicate {
				result.errors = append(result.errors, ValidationError{
					Field: field, Code: "duplicate", Message: "dependency is listed more than once",
				})
				continue
			}
			seenDependencies[dependencyID] = struct{}{}

			switch {
			case dependencyID == job.ID:
				result.errors = append(result.errors, ValidationError{
					Field: field, Code: "self_dependency", Message: "job cannot depend on itself",
				})
				result.graphValid = false
			case !containsJobID(jobIDs, dependencyID):
				result.errors = append(result.errors, ValidationError{
					Field: field, Code: "unknown_dependency",
					Message: fmt.Sprintf("dependency %q does not exist", dependencyID),
				})
				result.graphValid = false
			}
		}
	}

	return result
}

func containsJobID(jobIDs map[string]struct{}, id string) bool {
	_, exists := jobIDs[id]
	return exists
}
