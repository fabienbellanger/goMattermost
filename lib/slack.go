package mattermost

import (
	"fmt"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fatih/color"
)

// sendNotificationToSlack : Envoi sur le Webhook de Slack
func sendNotificationToSlack(payloadJSONEncoded []byte) {
	fmt.Println("Sending notification to Slack...")

	// Récupération des paramètres
	// ---------------------------
	hookURL = config.SlackHookURL
	hookPayload = config.SlackHookPayload

	// Envoi de la requête
	// -------------------
	response := sendNotificationToApplication(hookURL, hookPayload, payloadJSONEncoded)

	fmt.Print(" -> Slack response: ")
	color.Green(response.Status + "\n\n")
}
