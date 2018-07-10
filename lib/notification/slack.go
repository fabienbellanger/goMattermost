package notification

import (
	"fmt"
	"net/http"

	"github.com/fabienbellanger/goMattermost/config"
)

// sendNotificationToSlack : Envoi sur le Webhook de Slack
func sendNotificationToSlack(payloadJSONEncoded []byte, response chan<- *http.Response) {
	fmt.Println("Sending notification to Slack...")

	// Récupération des paramètres
	// ---------------------------
	hookURL = config.SlackHookURL
	hookPayload = config.SlackHookPayload

	// Envoi de la requête
	// -------------------
	response <- sendNotificationToApplication(hookURL, hookPayload, payloadJSONEncoded)
}
