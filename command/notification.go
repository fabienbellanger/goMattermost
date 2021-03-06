package command

import (
	"errors"
	"strings"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/database"
	"github.com/fabienbellanger/goMattermost/lib/notification"
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

// path 		: chemin vers le dépôt git
// repository 	: Nom du dépôt git
// branch 		: Branche de production (par défaut: master)
// applications : Liste des applications à notifier
var path, repository, branch, applications string

func init() {
	// Flags
	// -----
	NotificationCommand.Flags().StringVarP(&path, "path", "p", "", "Path")
	NotificationCommand.Flags().StringVarP(&repository, "repository", "r", "", "Repository")
	NotificationCommand.Flags().StringVarP(&branch, "branch", "b", "master", "Branch use for production")
	NotificationCommand.Flags().StringVarP(&applications, "applications", "a", "", "Applications list separate by commat withour whitespace (Ex: -a slack,mattermost)")
	NotificationCommand.Flags().BoolVarP(&NoDatabase, "no-database", "d", false, "Save data to database")
	NotificationCommand.MarkFlagRequired("path")
	NotificationCommand.MarkFlagRequired("repository")

	// Ajout de la commande à la commande racine
	rootCommand.AddCommand(NotificationCommand)
}

// NotificationCommand : Notification command
var NotificationCommand = &cobra.Command{
	Use:   "notification",
	Short: "Send message to Mattermost and/or Slack",

	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow(`

|-----------------------------------------------------------------------------------|
|                                                                                   |
| Send data of the last commit done on production branch to Mattermost and/or Slack |
|                                                                                   |
|-----------------------------------------------------------------------------------|

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
		notification.Launch(path, repository, branch, NoDatabase, SendToMattermost, SendToSlack)
	},
}
