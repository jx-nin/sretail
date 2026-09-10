package router_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jx-nin/sretail/services/catalog/internal/product"
	"github.com/jx-nin/sretail/services/catalog/internal/router"
)

func TestCatalogEndpoints(t *testing.T) {
	originalMode := gin.Mode()
	originalWriter := gin.DefaultWriter
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard
	t.Cleanup(func() {
		gin.SetMode(originalMode)
		gin.DefaultWriter = originalWriter
	})

	store := product.NewMemoryStore([]product.Product{
		{
			ID:         "prod-1",
			Name:       "Mechanical Keyboard",
			PriceCents: 12999,
		},
		{
			ID:         "prod-2",
			Name:       "Ergonomic Mouse",
			PriceCents: 10999,
		},
	})
	r := router.New(store)

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "lists products",
			path:       "/products",
			wantStatus: http.StatusOK,
			wantBody:   `[{"id":"prod-1","name":"Mechanical Keyboard","price_cents":12999},{"id":"prod-2","name":"Ergonomic Mouse","price_cents":10999}]`,
		},
		{
			name:       "returns a product by ID",
			path:       "/products/prod-1",
			wantStatus: http.StatusOK,
			wantBody:   `{"id":"prod-1","name":"Mechanical Keyboard","price_cents":12999}`,
		},
		{
			name:       "returns not found for an unknown product",
			path:       "/products/missing",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"product not found"}`,
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

			if got := recorder.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
				t.Errorf("Content-Type = %q, want %q", got, "application/json; charset=utf-8")
			}

			if got := recorder.Body.String(); got != tt.wantBody {
				t.Errorf("response body = %q, want %q", got, tt.wantBody)
			}
		})
	}
}
