package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func AssertResponse(t *testing.T, got *http.Response, wantStatusCode int, wantb []byte) {
	t.Helper()

	t.Cleanup(func() {
		_ = got.Body.Close()
	})
	gotb, err := io.ReadAll(got.Body)
	require.NoError(t, err, "reading response body")

	require.Equal(t, wantStatusCode, got.StatusCode, "validate status code")
	if len(gotb) == 0 && len(wantb) == 0 {
		return
	}

	AssertJSON(t, wantb, gotb)
}

func AssertJSON(t *testing.T, want, got []byte) {
	t.Helper()

	var jwant, jgot any
	err := json.Unmarshal(want, &jwant)
	require.NoError(t, err,
		"assert json: unmarshal want (bytes) to json struct")
	err = json.Unmarshal(got, &jgot)
	require.NoError(t, err,
		"assert json: unmarchal got (bytes) to json struct")
	diff := cmp.Diff(jwant, jgot)
	require.Equal(t, "", diff,
		"assert json: validate if json want and got have no diff")
}
