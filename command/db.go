package command

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry/bytefmt"
	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/database"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// Init : Initialisation de la base de données
var Init bool

// Dump : Dump de la base de données
var Dump bool

func init() {
	DbCommand.Flags().BoolVarP(&Init, "init", "i", false, "Database initialization")
	DbCommand.Flags().BoolVarP(&Dump, "dump", "d", false, "Database dump")

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

		if Init {
			// Initialisation
			var confirm string

			fmt.Println("If a database already exists, data will be deleted.")
			fmt.Print("Do you really want to initalize database (Y/n): ")
			_, err := fmt.Scanf("%s", &confirm)
			toolbox.CheckError(err, 1)

			if confirm == "Y" {
				fmt.Print("\n\n -> Database initialization:\t")

				database.InitDatabase()

				color.Green("Success\n\n")
			} else {
				fmt.Print("\n\n -> Database initialization:\t")
				color.Yellow("Operation aborded\n\n")
			}
		} else if Dump {
			// Dump
			fmt.Print("\n\n -> Database dump:\t")

			fileName, fileSize := database.DumpDatabase()

			color.Green(fileName + " (" + bytefmt.ByteSize(uint64(fileSize)) + ") successfully created\n\n")
		}
	},
}
