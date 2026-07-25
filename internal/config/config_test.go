package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	wantPort := 3333
	t.Setenv("PORT", fmt.Sprint(wantPort))

	got, err := New()
	require.NoError(t, err, "cannot create config")

	require.Equal(t, wantPort, got.Port, "config port")

	wantEnv := "dev"
	require.Equal(t, wantEnv, got.Env, "config env var")
}
