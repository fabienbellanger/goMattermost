package command

import (
	"errors"
	"fmt"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/database"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Force : Force l'initialisation de la base de données
var Force bool

func init() {
	DbCommand.Flags().BoolVarP(&Force, "force", "f", false, "Force database initialization")
	DbCommand.MarkFlagRequired("force")

	// Ajout de la commande à la commande racine
	rootCommand.AddCommand(DbCommand)
}

// DbCommand : Database command
var DbCommand = &cobra.Command{
	Use:   "db",
	Short: "Initialisation de la base de données",

	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow(`

|-------------------------|
|                         |
| Database initialization |
|                         |
|-------------------------|

		`)

		// Connexion à MySQL
		// -----------------
		if !config.IsDatabaseConfigCorrect() {
			err := errors.New("No or missing database information in settings file")
			toolbox.CheckError(err, 1)
		}

		database.Open()
		defer database.DB.Close()

		// Initialisation
		fmt.Print(" -> Database initialization: \t")

		if Force {
			database.InitDatabase()

			color.Green("Success\n\n")
		}
	},
}
