package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	tests := []struct {
		offset       int64
		limit        int64
		expectedFile string
	}{
		{
			offset:       0,
			limit:        0,
			expectedFile: "out_offset0_limit0",
		},
		{
			offset:       0,
			limit:        10,
			expectedFile: "out_offset0_limit10",
		},
		{
			offset:       0,
			limit:        1000,
			expectedFile: "out_offset0_limit1000",
		},
		{
			offset:       0,
			limit:        10000,
			expectedFile: "out_offset0_limit10000",
		},
		{
			offset:       100,
			limit:        1000,
			expectedFile: "out_offset100_limit1000",
		},
		{
			offset:       6000,
			limit:        1000,
			expectedFile: "out_offset6000_limit1000",
		},
	}

	for _, val := range tests {
		t.Run(val.expectedFile, func(t *testing.T) {
			temp, err := os.CreateTemp(os.TempDir(), "test_copy")
			require.NoError(t, err)

			err = Copy("testdata/input.txt", temp.Name(), val.offset, val.limit)
			require.NoError(t, err)

			expectedFile, err := os.OpenFile("testdata/"+val.expectedFile+".txt", os.O_RDONLY, 0644)
			require.NoError(t, err)

			sc1 := bufio.NewScanner(expectedFile)
			sc2 := bufio.NewScanner(temp)

			for {
				sc1Bool := sc1.Scan()
				sc2Bool := sc2.Scan()
				if !sc1Bool && !sc2Bool {
					break
				}

				require.Equal(t, sc1.Text(), sc2.Text())
			}

			err = expectedFile.Close()
			require.NoError(t, err)

			err = temp.Close()
			require.NoError(t, err)
			err = os.Remove(temp.Name())
			require.NoError(t, err)
		})
	}
}
