package notification

import (
	"fmt"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fatih/color"
)

// sendNotificationToMattermost : Envoi sur le Webhook de Mattermost
func sendNotificationToMattermost(payloadJSONEncoded []byte) {
	fmt.Println("Sending notification to Mattermost...")

	// Récupération des paramètres
	// ---------------------------
	hookURL = config.MattermostHookURL
	hookPayload = config.MattermostHookPayload

	// Envoi de la requête
	// -------------------
	response := sendNotificationToApplication(hookURL, hookPayload, payloadJSONEncoded)

	fmt.Print(" -> Mattermost response: ")
	color.Green(response.Status + "\n\n")
}
