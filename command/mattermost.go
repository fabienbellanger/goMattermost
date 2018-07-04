package command

import (
	"errors"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/database"
	"github.com/fabienbellanger/goMattermost/lib"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NoDatabase : Aucune opération ne sera faite en base de données
var NoDatabase bool
var path, repository string

func init() {
	// Flags
	// -----
	MattermostCommand.Flags().StringVarP(&path, "path", "p", "", "Path")
	MattermostCommand.Flags().StringVarP(&repository, "repository", "r", "", "Repository")
	MattermostCommand.MarkFlagRequired("path")
	MattermostCommand.MarkFlagRequired("repository")

	rootCommand.PersistentFlags().BoolVarP(&NoDatabase, "no-database", "d", false, "Save data to database")

	// Ajout de la commande à la commande racine
	rootCommand.AddCommand(MattermostCommand)
}

// MattermostCommand : Mattermost command
var MattermostCommand = &cobra.Command{
	Use:   "mattermost",
	Short: "Send message to Mattermost",

	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow(`
|------------------------------------------------------------------|
|                                                                  |
| Envoi des données du dernier commit fait sur master à Mattermost |
|                                                                  |
|------------------------------------------------------------------|

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

		// Envoi à Mattermost
		// ------------------
		if !config.IsMattermostConfigCorrect() {
			err := errors.New("No or missing Mattermost information in settings file")
			toolbox.CheckError(err, 1)
		}
		mattermost.Launch(path, repository, NoDatabase)
	},
}
