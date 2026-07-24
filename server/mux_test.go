package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/k20ku/see/testutil"
)

func TestNewMux(t *testing.T) {
	// ResponseWriter
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	sut := NewMux()
	sut.ServeHTTP(w, r)
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	want := []byte(`{"status" : "ok"}`)
	testutil.AssertResponse(t, resp, http.StatusOK, want)
}
