package command

import (
	"os"

	settings "github.com/fabienbellanger/goMattermost/config"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:     "goMattermost",
	Short:   "goMattermost send notification to Mattermost",
	Long:    "goMattermost send notification to Mattermost",
	Version: settings.Version,
}

// Execute starts Cobra
func Execute() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
