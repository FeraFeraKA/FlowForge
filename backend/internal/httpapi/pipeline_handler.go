package httpapi

import (
	"context"
	"net/http"

	"flowforge/backend/internal/pipeline"
)

type PipelineService interface {
	Create(ctx context.Context, params pipeline.CreateParams) (pipeline.Pipeline, error)
	GetByID(ctx context.Context, id string) (pipeline.Pipeline, error)
	List(ctx context.Context) ([]pipeline.Pipeline, error)
}

type pipelineHandler struct {
	service PipelineService
}

func (h pipelineHandler) create(w http.ResponseWriter, r *http.Request) {
	var request createPipelineRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	params, err := request.params()
	if err != nil {
		writeDomainError(w, err)
		return
	}
	created, err := h.service.Create(r.Context(), params)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	w.Header().Set("Location", "/api/pipelines/"+created.ID)
	writeJSON(w, http.StatusCreated, toPipelineResponse(created))
}

func (h pipelineHandler) list(w http.ResponseWriter, r *http.Request) {
	pipelines, err := h.service.List(r.Context())
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toPipelineListResponse(pipelines))
}

func (h pipelineHandler) get(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetByID(r.Context(), r.PathValue("pipelineID"))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toPipelineResponse(result))
}
