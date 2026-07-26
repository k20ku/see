package testutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	require.NoErrorf(t, err, "cannot read from file %q", path)

	return bt
}
