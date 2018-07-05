package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration type
type Configuration struct {
	Version  string
	Database struct {
		Driver   string
		Name     string
		User     string
		Password string
	}
	Mattermost struct {
		HookURL     string
		HookPayload string
	}
	Slack struct {
		HookURL string
	}
}

// Version of current application
var Version string

// ============================================================================
//
// Database
//
// ============================================================================

// DatabaseDriver : Driver
var DatabaseDriver string

// DatabaseName : Nom de la base de données
var DatabaseName string

// DatabaseUser : Utilisateur
var DatabaseUser string

// DatabasePassword : Mot de passe
var DatabasePassword string

// ============================================================================
//
// Mattermost
//
// ============================================================================

// MattermostHookURL : URL pour l'envoi de message
var MattermostHookURL string

// MattermostHookPayload : Payload pour l'envoi de message
var MattermostHookPayload string

// Init : Lecture du fichier de configuration
func Init() {
	// Lecture du fichier de configuration
	// -----------------------------------
	file, _ := os.Open("settings.json")
	defer file.Close()

	// Décodage du JSON
	// ----------------
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)

	if err != nil {
		fmt.Println("error:", err)
	}

	// Initialisation des variables
	// ----------------------------
	Version = configuration.Version

	DatabaseDriver = configuration.Database.Driver
	DatabaseName = configuration.Database.Name
	DatabaseUser = configuration.Database.User
	DatabasePassword = configuration.Database.Password

	MattermostHookURL = configuration.Mattermost.HookURL
	MattermostHookPayload = configuration.Mattermost.HookPayload
}

// IsMattermostConfigCorrect : La configuration de Mattermost est-elle correcte ?
func IsMattermostConfigCorrect() bool {
	return MattermostHookURL != "" && MattermostHookPayload != ""
}

// IsDatabaseConfigCorrect : La configuration de la base de données est-elle correcte ?
func IsDatabaseConfigCorrect() bool {
	return (DatabaseDriver != "" && DatabaseName != "" && DatabaseUser != "")
}
