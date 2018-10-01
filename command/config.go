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

		displayStatus()
	},
}

// convertBoolStatusToString : Converti un booléan en string
func convertBoolStatusToString(status bool) string {
	if status {
		return "[OK]"
	}

	return "[KO]"
}

// displayStatus : Affiche les status des paramètres
func displayStatus() {
	color.Yellow("\n\nSTATUS")
	color.Yellow("------\n")

	fmt.Print("    " + toolbox.Ucfirst("database") + "\t")
	if config.IsDatabaseConfigCorrect() {
		color.Green(convertBoolStatusToString(config.IsDatabaseConfigCorrect()))
	} else {
		color.Red(convertBoolStatusToString(config.IsDatabaseConfigCorrect()))
	}

	fmt.Print("    " + strings.ToUpper("jwt") + "\t\t")
	color.Green(convertBoolStatusToString(true))

	fmt.Print("    " + toolbox.Ucfirst("slack") + "\t")
	if config.IsSlackConfigCorrect() {
		color.Green(convertBoolStatusToString(config.IsSlackConfigCorrect()))
	} else {
		color.Red(convertBoolStatusToString(config.IsSlackConfigCorrect()))
	}

	fmt.Print("    " + toolbox.Ucfirst("mattermost") + "\t")
	if config.IsMattermostConfigCorrect() {
		color.Green(convertBoolStatusToString(config.IsMattermostConfigCorrect()))
	} else {
		color.Red(convertBoolStatusToString(config.IsMattermostConfigCorrect()))
	}

	fmt.Print("    " + strings.ToUpper("smtp") + "\t")
	if config.IsSMTPServerConfigValid() {
		color.Green(convertBoolStatusToString(config.IsSMTPServerConfigValid()))
	} else {
		color.Red(convertBoolStatusToString(config.IsSMTPServerConfigValid()))
	}

	fmt.Println()
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

	switch name {
	case "database":
		configuration = config.GetDatabaseConfig()
		break
	case "jwt":
		configuration = config.GetJWTConfig()
		break
	case "slack":
		configuration = config.GetSlackConfig()
		break
	case "mattermost":
		configuration = config.GetMattermostConfig()
		break
	case "smtp":
		configuration = config.GetSMTPServerConfig()
		break
	default:
		configuration = nil
	}

	if configuration != nil {
		color.Yellow("\n\n" + strings.ToUpper(name))
		// color.Yellow("===================\n")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		for key, value := range configuration {
			table.Append([]string{toolbox.Ucfirst(key), value})
		}

		table.Render()
	}
}
