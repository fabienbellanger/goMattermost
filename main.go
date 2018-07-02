package main

import (
	"github.com/fabienbellanger/goMattermost/command"
	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/database"
)

func main() {
	// Initialisation de la configuration
	// ----------------------------------
	config.Init()

	// Connexion à MySQL
	// -----------------
	database.Open()
	defer database.DB.Close()

	// Lancement de Cobra
	command.Execute()
}
