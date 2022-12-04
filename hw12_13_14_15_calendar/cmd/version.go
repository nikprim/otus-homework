package cmd

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "version app",
		Run: func(cmd *cobra.Command, args []string) {
			if err := json.NewEncoder(os.Stdout).Encode(struct {
				Release   string
				BuildDate string
				GitHash   string
			}{
				Release:   release,
				BuildDate: buildDate,
				GitHash:   gitHash,
			}); err != nil {
				log.Error().Err(err).Msg("error while decode version info")
			}
		},
	}
}
