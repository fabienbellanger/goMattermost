package command

import (
	"os"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:     "goMattermost",
	Short:   "goMattermost send notification to Mattermost",
	Long:    "goMattermost send notification to Mattermost",
	Version: config.Version,
}

// Execute starts Cobra
func Execute() {
	// Initialisation de la configuration
	// ----------------------------------
	config.Init()

	// Lancement de la commande racine
	// -------------------------------
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
