package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/fabienbellanger/goMattermost/toolbox"
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
	JWT struct {
		SecretKey string
		Exp       time.Duration
	}
	Mattermost struct {
		HookURL     string
		HookPayload string
	}
	Slack struct {
		HookURL     string
		HookPayload string
	}
	SMTP struct {
		Host     string
		Port     string
		Username string
		Password string
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
// JWT
//
// ============================================================================

// JWTSecretKey : Clé secrète JWT
var JWTSecretKey string

// JWTExp : Durée de validité du token
var JWTExp time.Duration

// ============================================================================
//
// Mattermost
//
// ============================================================================

// MattermostHookURL : URL pour l'envoi de message
var MattermostHookURL string

// MattermostHookPayload : Payload pour l'envoi de message
var MattermostHookPayload string

// ============================================================================
//
// Slack
//
// ============================================================================

// SlackHookURL : URL pour l'envoi de message
var SlackHookURL string

// SlackHookPayload : Payload pour l'envoi de message
var SlackHookPayload string

// ============================================================================
//
// SMTP
//
// ============================================================================

// SMTPHost : Host du serveur SMTP
var SMTPHost string

// SMTPPort : Port  du serveur SMTP
var SMTPPort string

// SMTPUsername : Username du serveur SMTP
var SMTPUsername string

// SMTPPassword : Password du serveur SMTP
var SMTPPassword string

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
	toolbox.CheckError(err, 0)

	// Initialisation des variables
	// ----------------------------
	Version = configuration.Version

	DatabaseDriver = configuration.Database.Driver
	DatabaseName = configuration.Database.Name
	DatabaseUser = configuration.Database.User
	DatabasePassword = configuration.Database.Password

	JWTSecretKey = configuration.JWT.SecretKey
	JWTExp = time.Hour * configuration.JWT.Exp

	MattermostHookURL = configuration.Mattermost.HookURL
	MattermostHookPayload = configuration.Mattermost.HookPayload

	SlackHookURL = configuration.Slack.HookURL
	SlackHookPayload = configuration.Slack.HookPayload

	SMTPHost = configuration.SMTP.Host
	SMTPPort = configuration.SMTP.Port
	SMTPUsername = configuration.SMTP.Username
	SMTPPassword = configuration.SMTP.Password
}

// IsMattermostConfigCorrect : La configuration de Mattermost est-elle correcte ?
func IsMattermostConfigCorrect() bool {
	return MattermostHookURL != "" && MattermostHookPayload != ""
}

// IsSlackConfigCorrect : La configuration de Slack est-elle correcte ?
func IsSlackConfigCorrect() bool {
	return SlackHookURL != "" && SlackHookPayload != ""
}

// IsDatabaseConfigCorrect : La configuration de la base de données est-elle correcte ?
func IsDatabaseConfigCorrect() bool {
	return (DatabaseDriver != "" && DatabaseName != "" && DatabaseUser != "")
}

// IsSMTPServerConfigValid : La configuration du serveur SMTP est-elle correcte ?
func IsSMTPServerConfigValid() bool {
	return (SMTPHost != "")
}
