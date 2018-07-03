package command

import (
	"github.com/fabienbellanger/goMattermost/lib"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var path, repository string

func init() {
	// Flags
	// -----
	MattermostCommand.Flags().StringVarP(&path, "path", "p", "", "Path")
	MattermostCommand.Flags().StringVarP(&repository, "repository", "r", "", "Repository")
	MattermostCommand.MarkFlagRequired("path")
	MattermostCommand.MarkFlagRequired("repository")

	// Ajout de la commande à la commande racine
	rootCommand.AddCommand(MattermostCommand)
}

// MattermostCommand : Web command
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

		// Envoi à Mattermost
		// ------------------
		mattermost.Launch(path, repository)
	},
}
