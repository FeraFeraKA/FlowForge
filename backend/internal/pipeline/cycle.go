package pipeline

type visitState uint8

const (
	visitStateUnvisited visitState = iota
	visitStateVisiting
	visitStateVisited
)

func hasDependencyCycle(jobs []JobDefinition) bool {
	dependencies := make(map[string][]string, len(jobs))
	for _, job := range jobs {
		dependencies[job.ID] = job.DependsOn
	}

	states := make(map[string]visitState, len(jobs))
	var visit func(string) bool
	visit = func(jobID string) bool {
		switch states[jobID] {
		case visitStateVisiting:
			return true
		case visitStateVisited:
			return false
		}

		states[jobID] = visitStateVisiting
		for _, dependencyID := range dependencies[jobID] {
			if visit(dependencyID) {
				return true
			}
		}
		states[jobID] = visitStateVisited
		return false
	}

	for jobID := range dependencies {
		if states[jobID] == visitStateUnvisited && visit(jobID) {
			return true
		}
	}

	return false
}
