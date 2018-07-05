package mattermost

import (
	"fmt"
)

// sendNotificationToMattermost : Envoi sur le Webhook de Slack
func sendNotificationToSlack(payloadJSONEncoded []byte) {
	fmt.Println("Sending notification to Slack...")

	// Construction de la requête
	// --------------------------
	/*data := url.Values{}
	data.Set("payload", string(payloadJSONEncoded))

	u, err := url.ParseRequestURI(hookURL + hookPayload)
	toolbox.CheckError(err, 4)
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	// Envoi de la requête
	// -------------------
	response, err := client.Do(r)
	toolbox.CheckError(err, 5)

	fmt.Print(" -> Slack response: ")
	color.Green(response.Status + "\n\n")*/
}
