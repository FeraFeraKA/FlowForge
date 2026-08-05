package httpapi

import (
	"context"
	"net/http"

	runmodel "flowforge/backend/internal/run"
)

type RunService interface {
	StartRun(ctx context.Context, pipelineID string) (runmodel.Run, error)
	GetByID(ctx context.Context, id string) (runmodel.Run, error)
	Cancel(ctx context.Context, id string) error
	Retry(ctx context.Context, id string) (runmodel.Run, error)
}

type runHandler struct {
	service RunService
}

func (h runHandler) start(w http.ResponseWriter, r *http.Request) {
	execution, err := h.service.StartRun(r.Context(), r.PathValue("pipelineID"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	w.Header().Set("Location", "/api/runs/"+execution.ID)
	writeJSON(w, http.StatusAccepted, toRunResponse(execution))
}

func (h runHandler) get(w http.ResponseWriter, r *http.Request) {
	execution, err := h.service.GetByID(r.Context(), r.PathValue("runID"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toRunResponse(execution))
}

func (h runHandler) cancel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("runID")
	if err := h.service.Cancel(r.Context(), id); err != nil {
		writeDomainError(w, err)
		return
	}
	execution, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, toRunResponse(execution))
}

func (h runHandler) retry(w http.ResponseWriter, r *http.Request) {
	execution, err := h.service.Retry(r.Context(), r.PathValue("runID"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	w.Header().Set("Location", "/api/runs/"+execution.ID)
	writeJSON(w, http.StatusAccepted, toRunResponse(execution))
}
