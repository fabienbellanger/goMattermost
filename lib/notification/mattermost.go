package notification

import (
	"fmt"
	"net/http"

	"github.com/fabienbellanger/goMattermost/config"
)

// sendNotificationToMattermost : Envoi sur le Webhook de Mattermost
func sendNotificationToMattermost(payloadJSONEncoded []byte, response chan<- *http.Response) {
	fmt.Println("Sending notification to Mattermost...")

	// Récupération des paramètres
	// ---------------------------
	hookURL = config.MattermostHookURL
	hookPayload = config.MattermostHookPayload

	// Envoi de la requête
	// -------------------
	response <- sendNotificationToApplication(hookURL, hookPayload, payloadJSONEncoded)
}
