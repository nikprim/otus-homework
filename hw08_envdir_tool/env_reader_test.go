package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	want := Environment{
		"BAR":   EnvValue{Value: "bar"},
		"EMPTY": EnvValue{NeedRemove: true},
		"FOO":   EnvValue{Value: "   foo\nwith new line"},
		"HELLO": EnvValue{Value: "\"hello\""},
		"UNSET": EnvValue{NeedRemove: true},
	}

	envs, err := ReadDir("./testdata/env")

	require.NoError(t, err)
	require.Equal(t, want, envs)
}

func TestEnvironmentToStrings(t *testing.T) {
	envs := Environment{
		"": {
			Value:      "",
			NeedRemove: false,
		},
		"FOO": {
			Value:      "foo",
			NeedRemove: false,
		},
		"BAR": {
			Value:      "bar",
			NeedRemove: true,
		},
		"EMPTY": {
			Value:      "",
			NeedRemove: false,
		},
	}

	want := []string{
		"FOO=foo",
		"BAR=",
		"EMPTY=",
	}

	require.Equal(t, want, envs.toStrings())
}

func TestGetValueFromLine(t *testing.T) {
	tests := []struct {
		line     string
		expected string
	}{
		{
			line:     "foo",
			expected: "foo",
		},
		{
			line:     "foo  ",
			expected: "foo",
		},
		{
			line:     "foo\t\t  \t \t",
			expected: "foo",
		},
		{
			line:     "foo" + string([]byte{0x00}) + "newline\t\t   ",
			expected: "foo\nnewline",
		},
	}

	for _, test := range tests {
		t.Run(test.line, func(t *testing.T) {
			result := getValueFromLine(test.line)
			require.Equal(t, test.expected, result)
		})
	}
}
