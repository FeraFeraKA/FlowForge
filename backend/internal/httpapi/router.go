package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

type Dependencies struct {
	Pipelines      PipelineService
	Runs           RunService
	AllowedOrigins []string
}

func NewRouter(dependencies Dependencies) http.Handler {
	mux := http.NewServeMux()
	pipelines := pipelineHandler{service: dependencies.Pipelines}
	runs := runHandler{service: dependencies.Runs}
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /api/pipelines", pipelines.create)
	mux.HandleFunc("GET /api/pipelines", pipelines.list)
	mux.HandleFunc("GET /api/pipelines/{pipelineID}", pipelines.get)
	mux.HandleFunc("POST /api/pipelines/{pipelineID}/runs", runs.start)
	mux.HandleFunc("GET /api/runs/{runID}", runs.get)
	mux.HandleFunc("POST /api/runs/{runID}/cancel", runs.cancel)
	mux.HandleFunc("POST /api/runs/{runID}/retry", runs.retry)
	return corsMiddleware(dependencies.AllowedOrigins, mux)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(healthResponse{Status: "ok"}); err != nil {
		slog.Error("write health response", "error", err)
	}
}
