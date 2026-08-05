package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"flowforge/backend/internal/pipeline"
)

func TestCreatePipelineEndpoint(t *testing.T) {
	service := pipeline.NewService(pipeline.NewMemoryRepository())
	router := NewRouter(Dependencies{Pipelines: service})
	body := `{
		"name":"checks","maxParallel":2,"jobs":[{
			"id":"tests","name":"Tests","type":"delay",
			"dependsOn":[],"maxRetries":1,"retryDelay":"10ms",
			"config":{"duration":"20ms"}
		}]}`
	request := httptest.NewRequest(http.MethodPost, "/api/pipelines", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusCreated, response.Body.String())
	}
	if location := response.Header().Get("Location"); !strings.HasPrefix(location, "/api/pipelines/") {
		t.Fatalf("Location = %q, want pipeline URL", location)
	}
	for _, fragment := range []string{`"name":"checks"`, `"maxParallel":2`, `"duration":"20ms"`} {
		if !strings.Contains(response.Body.String(), fragment) {
			t.Errorf("response does not contain %s: %s", fragment, response.Body.String())
		}
	}
}

func TestCreatePipelineValidationError(t *testing.T) {
	service := pipeline.NewService(pipeline.NewMemoryRepository())
	router := NewRouter(Dependencies{Pipelines: service})
	request := httptest.NewRequest(http.MethodPost, "/api/pipelines", strings.NewReader(
		`{"name":"broken","maxParallel":0,"jobs":[]}`,
	))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnprocessableEntity)
	}
	if !strings.Contains(response.Body.String(), `"code":"validation_failed"`) {
		t.Fatalf("body = %s, want validation error", response.Body.String())
	}
}
