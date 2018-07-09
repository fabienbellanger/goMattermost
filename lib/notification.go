package mattermost

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/fabienbellanger/goMattermost/models"
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
}

// Launch : Lancement du traitement
func Launch(path, repository string, noDatabase, sendToMattermost, sendToSlack bool) {
	// Récupération du dernier commit Git de master
	// --------------------------------------------
	gitLogOutput := retrieveCommit(path)

	// Formattage du commit et du repository
	// -------------------------------------
	commit := models.CommitInformation{}
	formatGitCommit(gitLogOutput, &commit)
	formatRepositoryName(&repository)

	// Formattage du payload
	// ---------------------
	payloadJSONEncodedMattermost := formatPayloadMattermost(repository, commit)
	payloadJSONEncodedSlack := formatPayloadSlack(repository, commit)

	// Envoi à Mattermost
	// ------------------
	if sendToMattermost {
		sendNotificationToMattermost(payloadJSONEncodedMattermost)
	}

	// Envoi à Slack
	// -------------
	if sendToSlack {
		sendNotificationToSlack(payloadJSONEncodedSlack)
	}

	// Enregistrement du commit en base de données
	// -------------------------------------------
	if !noDatabase {
		fmt.Println("Inserting commit into database...")

		commitDB, err := models.AddCommit(repository, commit)

		if err != nil {
			color.Red(" -> Error during inserting commit in database\n\n")
		} else {
			fmt.Print(" -> Commit inserted with ID: ")
			color.Green(strconv.FormatInt(int64(commitDB.ID), 10) + "\n\n")
		}
	}
}

// formatRepositoryName : Formattage du nom du répository
func formatRepositoryName(repository *string) {
	s := *repository
	lastIndex := strings.LastIndex(s, "/")

	repositoryLen := len(s)
	if lastIndex == repositoryLen-1 {
		*repository = s[0 : repositoryLen-1]
	}
}

// retrieveCommit : Récupération du dernier commit Git de master
func retrieveCommit(path string) []byte {
	gitLogCommand := exec.Command("git",
		"log",
		"-1",
		"--pretty=format:<%an>%n<%s>%n<%b>",
		"--first-parent", "master",
	)
	gitLogCommand.Dir = path
	gitLogOutput, err := gitLogCommand.Output()
	toolbox.CheckError(err, 1)

	return gitLogOutput
}

// formatGitCommit : Formattage du commit
func formatGitCommit(gitLogOutput []byte, commit *models.CommitInformation) {
	message := ""
	regex := regexp.MustCompile("(?m)(?s)<(.*)>\n<(.*)>\n<(.*)>")

	for _, match := range regex.FindAllSubmatch(gitLogOutput, -1) {
		if len(match) == 4 {
			commit.Author = string(match[1])
			commit.Subject = string(match[2])
			message = string(match[3])
		}
	}

	// Mise en forme du message
	// ------------------------
	if message != "" {
		regexMessage := regexp.MustCompile("(?s)Version : (.*)\nDev : (.*)\n(?:Test|Tests) : (.*)\nDescription :\n(.*)")

		for _, matchMessage := range regexMessage.FindAllSubmatch([]byte(message), -1) {
			if len(matchMessage) == 5 {
				commit.Version = string(matchMessage[1])
				commit.Developers = string(matchMessage[2])
				commit.Testers = string(matchMessage[3])
				commit.Description = string(matchMessage[4])
			}
		}
	}
}

// isCommitValid : Les informations du commit sont-elles valides ?
func isCommitValid(commit models.CommitInformation) bool {
	return (commit.Author != "" ||
		commit.Subject != "" ||
		commit.Version != "" ||
		commit.Developers != "" ||
		commit.Testers != "" ||
		commit.Description != "")
}

// formatPayloadMattermost : Mise en forme du payload au format Markdown
func formatPayloadMattermost(repository string, commit models.CommitInformation) []byte {
	if !isCommitValid(commit) {
		err := errors.New("No Git repository found")
		toolbox.CheckError(err, 2)
	}

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
func formatPayloadSlack(repository string, commit models.CommitInformation) []byte {
	if !isCommitValid(commit) {
		err := errors.New("No Git repository found")
		toolbox.CheckError(err, 2)
	}

	// Création du payload à transmettre
	// ---------------------------------
	payload := payload{
		Text:      "",
		IconEmoji: ":ghost:",
		Channel:   "#mep",
		Username:  "mep " + repository,
	}
	payload.Text = "Mise en production\n"
	payload.Text += " " + toolbox.Ucfirst(repository)

	if commit.Version != "" {
		payload.Text += " - v" + commit.Version
	}
	payload.Text += "\n"

	if commit.Subject != "" {
		payload.Text += "Sujet : " + toolbox.Ucfirst(commit.Subject) + "\n"
	}

	if commit.Author != "" {
		payload.Text += "Auteur : " + toolbox.Ucfirst(commit.Author) + "\n"
	}

	if commit.Developers != "" {
		payload.Text += "Développeur(s) : " + commit.Developers + "\n"
	}

	if commit.Testers != "" {
		payload.Text += "Testeur(s) : " + commit.Testers + "\n"
	}

	if commit.Description != "" {
		payload.Text += "Description :\n" + commit.Description + "\n"
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
