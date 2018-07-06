package command

import (
	"errors"
	"strings"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/database"
	notification "github.com/fabienbellanger/goMattermost/lib"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NoDatabase : Aucune opération ne sera faite en base de données
var NoDatabase bool

// SendToMattermost : Envoi à Mattermost ?
var SendToMattermost bool

// SendToSlack : Envoi à Slack ?
var SendToSlack bool

var path, repository, applications string

func init() {
	// Flags
	// -----
	NotificationCommand.Flags().StringVarP(&path, "path", "p", "", "Path")
	NotificationCommand.Flags().StringVarP(&repository, "repository", "r", "", "Repository")
	NotificationCommand.Flags().StringVarP(&applications, "applications", "a", "", "Applications list separate by commat withour whitespace (Ex: -a slack,mattermost)")
	NotificationCommand.MarkFlagRequired("path")
	NotificationCommand.MarkFlagRequired("repository")

	rootCommand.PersistentFlags().BoolVarP(&NoDatabase, "no-database", "d", false, "Save data to database")

	// Ajout de la commande à la commande racine
	rootCommand.AddCommand(NotificationCommand)
}

// NotificationCommand : Notification command
var NotificationCommand = &cobra.Command{
	Use:   "mattermost",
	Short: "Send message to Mattermost and/or Slack",

	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow(`
|-------------------------------------------------------------------------------|
|                                                                               |
| Send data of the last commit done on master branch to Mattermost and/or Slack |
|                                                                               |
|-------------------------------------------------------------------------------|

		`)

		// Connexion à MySQL
		// -----------------
		if !NoDatabase {
			if !config.IsDatabaseConfigCorrect() {
				err := errors.New("No or missing database information in settings file")
				toolbox.CheckError(err, 1)
			}

			database.Open()
			defer database.DB.Close()
		}

		// Liste des applications à notifier
		// ---------------------------------
		applicationsArray := strings.Split(applications, ",")
		SendToMattermost, _ = toolbox.InArray("mattermost", applicationsArray)
		SendToSlack, _ = toolbox.InArray("slack", applicationsArray)

		// Configuration Mattermost OK ?
		// -----------------------------
		if SendToMattermost && !config.IsMattermostConfigCorrect() {
			err := errors.New("No or missing Mattermost information in settings file")
			toolbox.CheckError(err, 1)
		}

		// Configuration Slack OK ?
		// ------------------------
		if SendToSlack && !config.IsSlackConfigCorrect() {
			err := errors.New("No or missing Slack information in settings file")
			toolbox.CheckError(err, 1)
		}

		// Envoi de la notification
		// ------------------------
		notification.Launch(path, repository, NoDatabase, SendToMattermost, SendToSlack)
	},
}
