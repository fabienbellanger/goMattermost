package command

import (
	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/toolbox"
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
	err := rootCommand.Execute()
	toolbox.CheckError(err, 1)
}
