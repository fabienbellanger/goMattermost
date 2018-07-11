package notification

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/fabienbellanger/goMattermost/model"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
)

var hookURL, hookPayload string

// payload structure
type payload struct {
	Text      string `json:"text"`
	Username  string `json:"username"`
	Channel   string `json:"channel"`
	IconEmoji string `json:"icon_emoji"`
	Markdown  bool   `json:"mrkdwn"`
}

// Launch : Lancement du traitement
func Launch(path, repository string, noDatabase, sendToMattermost, sendToSlack bool) {
	// Récupération du dernier commit Git de master
	// --------------------------------------------
	gitLogOutput := model.RetrieveCommit(path)

	// Formattage du commit et du repository
	// -------------------------------------
	commit := model.CommitInformation{}
	model.FormatGitCommit(gitLogOutput, &commit)
	model.FormatRepositoryName(&repository)

	// Formattage du payload
	// ---------------------
	payloadJSONEncodedMattermost := formatPayloadMattermost(repository, commit)
	payloadJSONEncodedSlack := formatPayloadSlack(repository, commit)

	// Envoi à Mattermost
	// ------------------
	var responseMattermostChan chan *http.Response

	if sendToMattermost {
		responseMattermostChan = make(chan *http.Response)
		go sendNotificationToMattermost(payloadJSONEncodedMattermost, responseMattermostChan)
	}

	// Envoi à Slack
	// -------------
	var responseSlackChan chan *http.Response

	if sendToSlack {
		responseSlackChan = make(chan *http.Response)
		go sendNotificationToSlack(payloadJSONEncodedSlack, responseSlackChan)
	}

	// Enregistrement du commit en base de données
	// -------------------------------------------
	var commitDBChan chan model.CommitDB
	var errorInsertChan chan error
	if !noDatabase {
		fmt.Println("Inserting commit into database...")

		commitDBChan = make(chan model.CommitDB)
		errorInsertChan = make(chan error)

		go model.AddCommit(repository, commit, commitDBChan, errorInsertChan)
	}

	// Réponses asynchrones
	// --------------------
	if sendToMattermost {
		responseMattermost := <-responseMattermostChan

		fmt.Print(" -> Mattermost response: \t")
		color.Green(responseMattermost.Status + "\n")
	}

	if sendToSlack {
		responseSlack := <-responseSlackChan

		fmt.Print(" -> Slack response: \t\t")
		color.Green(responseSlack.Status + "\n")
	}

	if !noDatabase {
		if <-errorInsertChan != nil {
			color.Red(" -> Error during inserting commit in database\n")
		} else {
			commitDB := <-commitDBChan

			fmt.Print(" -> Commit inserted with ID: \t")
			color.Green(strconv.FormatInt(int64(commitDB.ID), 10) + "\n")
		}
	}

	fmt.Println("")
}

// formatPayloadMattermost : Mise en forme du payload au format Markdown
func formatPayloadMattermost(repository string, commit model.CommitInformation) []byte {
	if !model.IsCommitValid(commit) {
		err := errors.New("No Git repository found")
		toolbox.CheckError(err, 2)
	}

	// Date et heure de la mise en production
	// --------------------------------------
	datetime := time.Now().Format("02/01/2006 à 15:04")

	// Création du payload à transmettre
	// ---------------------------------
	payload := payload{
		Text: "",
	}
	payload.Text = "### Mise en production\n"
	payload.Text += "#### " + toolbox.Ucfirst(repository)

	if commit.Version != "" {
		payload.Text += " - v" + commit.Version
	}
	payload.Text += "\n"

	if commit.Subject != "" {
		payload.Text += "| Sujet |" + toolbox.Ucfirst(commit.Subject) + "|\n"
	}

	payload.Text += "|:---|:---|\n"

	payload.Text += "| Date et heure | " + datetime + " |\n"

	if commit.Author != "" {
		payload.Text += "| Auteur |" + toolbox.Ucfirst(commit.Author) + "|\n"
	}

	if commit.Developers != "" {
		payload.Text += "| Développeur(s) |" + commit.Developers + "|\n"
	}

	if commit.Testers != "" {
		payload.Text += "| Testeur(s) |" + commit.Testers + "|\n"
	}

	if commit.Description != "" {
		payload.Text += "#### Description :\n" + commit.Description + "\n"
	}

	payloadJSONEncoded, err := json.Marshal(payload)
	toolbox.CheckError(err, 3)

	return payloadJSONEncoded
}

// formatPayloadSlack : Mise en forme du payload au format Texte
func formatPayloadSlack(repository string, commit model.CommitInformation) []byte {
	if !model.IsCommitValid(commit) {
		err := errors.New("No Git repository found")
		toolbox.CheckError(err, 2)
	}

	// Date et heure de la mise en production
	// --------------------------------------
	datetime := time.Now().Format("02/01/2006 à *15:04*")

	// Création du payload à transmettre
	// ---------------------------------
	payload := payload{
		Text:      "",
		IconEmoji: ":ghost:",
		Channel:   "#mep",
		Username:  "Mise en production : " + toolbox.Ucfirst(repository),
		Markdown:  true,
	}

	if commit.Version != "" {
		payload.Username += " - v" + commit.Version
	}

	payload.Text += "_Date et heure_ : " + datetime + "\n"

	if commit.Subject != "" {
		payload.Text += "_Sujet_ : " + toolbox.Ucfirst(commit.Subject) + "\n"
	}

	if commit.Author != "" {
		payload.Text += "_Auteur_ : " + toolbox.Ucfirst(commit.Author) + "\n"
	}

	if commit.Developers != "" {
		payload.Text += "_Développeur(s)_ : " + commit.Developers + "\n"
	}

	if commit.Testers != "" {
		payload.Text += "_Testeur(s)_ : " + commit.Testers + "\n"
	}

	if commit.Description != "" {
		payload.Text += "_Description_ :\n" + commit.Description + "\n"
	}

	payloadJSONEncoded, err := json.Marshal(payload)
	toolbox.CheckError(err, 3)

	return payloadJSONEncoded
}

// sendNotificationToApplication : Envoi du webhook à l'applicatif
func sendNotificationToApplication(hookURL, hookPayload string, payloadJSONEncoded []byte) *http.Response {
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

	return response
}
