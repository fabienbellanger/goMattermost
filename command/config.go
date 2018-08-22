package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
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
		displayParameters("database")
		displayParameters("jwt")
		displayParameters("slack")
		displayParameters("mattermost")
		displayParameters("smtp")

		fmt.Println()
	},
}

// displayVersion : Affiche la version courante
func displayVersion() {
	version := config.Version

	fmt.Print("Version:\t")
	color.Green(version + "\n")
}

// displayParameters : Affiche les paramètres
func displayParameters(name string) {
	var configuration map[string]string
	var status bool

	switch name {
	case "database":
		configuration = config.GetDatabaseConfig()
		status = config.IsDatabaseConfigCorrect()
	case "jwt":
		configuration = config.GetJWTConfig()
		status = true
	case "slack":
		configuration = config.GetSlackConfig()
		status = config.IsSlackConfigCorrect()
	case "mattermost":
		configuration = config.GetMattermostConfig()
		status = config.IsMattermostConfigCorrect()
	case "smtp":
		configuration = config.GetSMTPServerConfig()
		status = config.IsSMTPServerConfigValid()
	default:
		configuration = nil
		status = false
	}

	if configuration != nil {
		color.Yellow("\n\n" + strings.ToUpper(name))
		color.Yellow("===================\n")

		// Status
		// ------
		fmt.Print("Status:\t")

		if status {
			color.Green("OK\n\n")
		} else {
			color.Red("KO\n\n")
		}

		// Configuration
		// -------------
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		for key, value := range configuration {
			table.Append([]string{toolbox.Ucfirst(key), value})
		}

		table.Render()
	}
}
