package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
	"flowforge/backend/internal/runmanager"
)

type errorEnvelope struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []errorDetail `json:"details,omitempty"`
}

type errorDetail struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		slog.Error("write JSON response", "error", err)
	}
}

func writeDomainError(w http.ResponseWriter, err error) {
	var validationErrors pipeline.ValidationErrors
	switch {
	case errors.As(err, &validationErrors):
		details := make([]errorDetail, len(validationErrors))
		for index, validationErr := range validationErrors {
			details[index] = errorDetail{
				Field: validationErr.Field, Code: validationErr.Code, Message: validationErr.Message,
			}
		}
		writeJSON(w, http.StatusUnprocessableEntity, errorEnvelope{Error: apiError{
			Code: "validation_failed", Message: "pipeline validation failed", Details: details,
		}})
	case errors.Is(err, pipeline.ErrPipelineNotFound), errors.Is(err, runmodel.ErrRunNotFound):
		writeAPIError(w, http.StatusNotFound, "not_found", "resource not found")
	case errors.Is(err, runmodel.ErrInvalidTransition),
		errors.Is(err, runmanager.ErrRunNotActive),
		errors.Is(err, runmanager.ErrRetryNotAllowed):
		writeAPIError(w, http.StatusConflict, "state_conflict", "operation conflicts with current state")
	default:
		slog.Error("unexpected HTTP handler error", "error", err)
		writeAPIError(w, http.StatusInternalServerError, "internal_error", "internal server error")
	}
}

func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorEnvelope{Error: apiError{Code: code, Message: message}})
}
