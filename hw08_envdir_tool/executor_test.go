package main

import (
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	envs, err := ReadDir("./testdata/env")
	require.NoError(t, err)

	err = os.Setenv("HELLO", "SHOULD_REPLACE")
	require.NoError(t, err)
	err = os.Setenv("FOO", "SHOULD_REPLACE")
	require.NoError(t, err)
	err = os.Setenv("UNSET", "SHOULD_REPLACE")
	require.NoError(t, err)
	err = os.Setenv("ADDED", "from original env")
	require.NoError(t, err)
	err = os.Setenv("EMPTY", "SHOULD_REPLACE")
	require.NoError(t, err)

	expected := `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2
`

	cmd := []string{
		"/bin/bash",
		"./testdata/echo.sh",
		"arg1=1",
		"arg2=2",
	}

	origin := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	code := RunCmd(cmd, envs)
	require.Equal(t, 0, code)

	err = w.Close()
	require.NoError(t, err)

	out, err := io.ReadAll(r)
	require.NoError(t, err)

	os.Stdout = origin

	require.Equal(t, expected, string(out))
}
