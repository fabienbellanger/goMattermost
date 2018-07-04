package command

import (
	"errors"
	"os"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/database"
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

	// Connexion à MySQL
	// -----------------
	if !NoDatabase {
		if len(config.DatabaseDriver) == 0 || len(config.DatabaseName) == 0 || len(config.DatabaseUser) == 0 {
			err := errors.New("No or missing database information in settings file")
			toolbox.CheckError(err, 1)
		}

		database.Open()
		defer database.DB.Close()
	}

	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
