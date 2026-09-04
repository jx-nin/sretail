package router_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jx-nin/sretail/services/gateway/internal/router"
)

func TestHealthEndpoints(t *testing.T) {
	originalMode := gin.Mode()
	originalWriter := gin.DefaultWriter
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard
	t.Cleanup(func() {
		gin.SetMode(originalMode)
		gin.DefaultWriter = originalWriter
	})

	r := router.New()

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantJSON   map[string]string
	}{
		{
			name:       "health returns service status",
			path:       "/health",
			wantStatus: http.StatusOK,
			wantJSON: map[string]string{
				"service": "gateway",
				"status":  "ok",
			},
		},
		{
			name:       "liveness reports success without a body",
			path:       "/health/live",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "readiness reports success without a body",
			path:       "/health/ready",
			wantStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)

			r.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status code = %d, want %d", recorder.Code, tt.wantStatus)
			}

			if tt.wantJSON == nil {
				if recorder.Body.Len() != 0 {
					t.Fatalf("body = %q, want an empty body", recorder.Body.String())
				}
				return
			}

			if got := recorder.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
				t.Errorf("Content-Type = %q, want %q", got, "application/json; charset=utf-8")
			}

			var gotJSON map[string]string
			if err := json.Unmarshal(recorder.Body.Bytes(), &gotJSON); err != nil {
				t.Fatalf("decode response body: %v", err)
			}

			if !reflect.DeepEqual(gotJSON, tt.wantJSON) {
				t.Errorf("response body = %#v, want %#v", gotJSON, tt.wantJSON)
			}
		})
	}
}
