package command

import (
	"errors"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/database"
	mail "github.com/fabienbellanger/goMattermost/lib"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	// Ajout de la commande à la commande racine
	rootCommand.AddCommand(MailCommand)
}

// MailCommand : Mail command
var MailCommand = &cobra.Command{
	Use:   "mail",
	Short: "Send mail with commits list of yesterday",

	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow(`

|------------------------------------------|
|                                          |
| Send mail with commits list of yesterday |
|                                          |
|------------------------------------------|

		`)

		if !config.IsDatabaseConfigCorrect() {
			err := errors.New("No or missing database information in settings file")
			toolbox.CheckError(err, 1)
		}

		database.Open()
		defer database.DB.Close()

		// Liste des applications à notifier
		// ---------------------------------
		mail.SendCommitsByMail()
	},
}
