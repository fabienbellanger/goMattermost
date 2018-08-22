package command

import (
	"fmt"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	// Ajout de la commande à la commande racine
	rootCommand.AddCommand(ConfigCommand)
}

// ConfigCommand : Database command
var ConfigCommand = &cobra.Command{
	Use:   "config",
	Short: "Vérification de la configuration",

	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow(`

|----------------|
|                |
| Settings check |
|                |
|----------------|

		`)

		displayVersion()
		displayDatabaseStatus()
		displayJWTStatus()
	},
}

// displayVersion : Affiche la version courante
func displayVersion() {
	version := config.Version

	fmt.Print("Version:\t")
	color.Green(version + "\n")
}

// displayDatabaseStatus : Affiche la configuration de la base de données
func displayDatabaseStatus() {
	configuration := config.GetDatabaseConfig()
	status := config.IsDatabaseConfigCorrect()

	color.Yellow("\n\nDATABASE")
	color.Yellow("========\n")

	// Status
	// ------
	fmt.Print("Status:\t\t")

	if status {
		color.Green("OK\n\n")
	} else {
		color.Red("KO\n\n")
	}

	// Configuration
	// -------------
	for key, value := range configuration {
		if key == "password" {
			fmt.Print(toolbox.Ucfirst(key) + ":\t")
		} else {
			fmt.Print(toolbox.Ucfirst(key) + ":\t\t")
		}
		color.Green(value + "\n")
	}
}

// displayJWTStatus : Affiche la configuration de JWT
func displayJWTStatus() {
	configuration := config.GetJWTConfig()
	status := true

	color.Yellow("\n\nJWT")
	color.Yellow("===\n")

	// Status
	// ------
	fmt.Print("Status:\t\t")

	if status {
		color.Green("OK\n\n")
	} else {
		color.Red("KO\n\n")
	}

	// Configuration
	// -------------
	for key, value := range configuration {
		if key == "secret key" {
			fmt.Print(toolbox.Ucfirst(key) + ":\t")
		} else {
			fmt.Print(toolbox.Ucfirst(key) + ":\t\t")
		}
		color.Green(value + "\n")
	}
}
