package mattermost

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
)

// sendToMattermost : Envoi sur le Webhook de Mattermost
func sendToMattermost(payloadJSONEncoded []byte) {
	// Construction de la requête
	// --------------------------
	data := url.Values{}
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

	fmt.Print("Mattermost response: ")
	color.Green(response.Status + "\n\n")
}
