package command

import (
	"errors"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/database"
	router "github.com/fabienbellanger/goMattermost/router"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var port, defaultPort int

func init() {
	// Flag
	// ----
	defaultPort = 1323
	WebCommand.Flags().IntVarP(&port, "port", "p", defaultPort, "listened port")

	// Ajout de la commande à la commande racine
	rootCommand.AddCommand(WebCommand)
}

// WebCommand : Web command
var WebCommand = &cobra.Command{
	Use:   "web",
	Short: "Launch the web server",

	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow(`

|--------------------------|
|                          |
| Lancement du serveur Web |
|                          |
|--------------------------|`)

		// Test du port
		// ------------
		if port < 1000 || port > 10000 {
			port = defaultPort
		}

		// Connexion à MySQL
		// -----------------
		if !config.IsDatabaseConfigCorrect() {
			err := errors.New("No or missing database information in settings file")
			toolbox.CheckError(err, 1)
		}

		database.Open()
		defer database.DB.Close()

		// Lancement du serveur web
		// ------------------------
		router.StartServer(port)
	},
}
