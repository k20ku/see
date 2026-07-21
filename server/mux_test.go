package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMux(t *testing.T) {
	// ResponseWriter
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	sut := NewMux()
	sut.ServeHTTP(w, r)
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	require.Equal(t, resp.StatusCode, http.StatusOK, "request status is ok?")
	got, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "reading response body")

	want := `{"status" : "ok"}`
	require.Equal(t, want, string(got), "validating response json")
}
