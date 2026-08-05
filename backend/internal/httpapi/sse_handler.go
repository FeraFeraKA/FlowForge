package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"flowforge/backend/internal/events"
)

const sseHeartbeatInterval = 15 * time.Second

type EventSource interface {
	Subscribe(runID string) (<-chan events.Event, func())
}

type sseHandler struct {
	runs   RunService
	events EventSource
}

func (h sseHandler) stream(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("runID")
	if _, err := h.runs.GetByID(r.Context(), runID); err != nil {
		writeDomainError(w, err)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeAPIError(w, http.StatusInternalServerError, "stream_unsupported", "streaming is not supported")
		return
	}

	stream, unsubscribe := h.events.Subscribe(runID)
	defer unsubscribe()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	heartbeat := time.NewTicker(sseHeartbeatInterval)
	defer heartbeat.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case event, open := <-stream:
			if !open || writeSSE(w, event) != nil {
				return
			}
			flusher.Flush()
		case <-heartbeat.C:
			if _, err := fmt.Fprint(w, ": heartbeat\n\n"); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

func writeSSE(w http.ResponseWriter, event events.Event) error {
	data, err := json.Marshal(toEventResponse(event))
	if err != nil {
		return fmt.Errorf("encode SSE event: %w", err)
	}
	if _, err := fmt.Fprintf(w, "id: %s\nevent: %s\ndata: %s\n\n", event.ID, event.Type, data); err != nil {
		return fmt.Errorf("write SSE event: %w", err)
	}
	return nil
}
