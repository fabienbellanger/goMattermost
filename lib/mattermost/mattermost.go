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

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
)

var hookURL, hookPayload string

// payload structure
type payload struct {
	Text string `json:"text"`
}

// commitInformation structure
type commitInformation struct {
	author      string
	subject     string
	version     string
	developers  string
	testers     string
	tickets     string
	description string
}

// Launch : Lancement du traitement
func Launch(path, repository string) {
	// Récupération du dernier commit Git de master
	// --------------------------------------------
	gitLogOutput := retrieveCommit(path)

	// Formattage du commit
	// --------------------
	commit := commitInformation{}
	formatGitCommit(gitLogOutput, &commit)

	// Envoi à Mattermost
	// ------------------
	sendToMattermost(repository, commit)
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
func formatGitCommit(gitLogOutput []byte, commit *commitInformation) {
	message := ""
	regex := regexp.MustCompile("(?m)(?s)<(.*)>\n<(.*)>\n<(.*)>")

	for _, match := range regex.FindAllSubmatch(gitLogOutput, -1) {
		if len(match) == 4 {
			commit.author = string(match[1])
			commit.subject = string(match[2])
			message = string(match[3])
		}
	}

	// Mise en forme du message
	// ------------------------
	if message != "" {
		regexMessage := regexp.MustCompile("(?s)Version : (.*)\nDev : (.*)\n(?:Test|Tests) : (.*)\nDescription :\n(.*)")

		for _, matchMessage := range regexMessage.FindAllSubmatch([]byte(message), -1) {
			if len(matchMessage) == 5 {
				commit.version = string(matchMessage[1])
				commit.developers = string(matchMessage[2])
				commit.testers = string(matchMessage[3])
				commit.description = string(matchMessage[4])
			}
		}
	}
}

// isCommitValid : Les informations du commit sont-elles valides ?
func isCommitValid(commit commitInformation) bool {
	return (commit.author != "" ||
		commit.subject != "" ||
		commit.version != "" ||
		commit.developers != "" ||
		commit.testers != "" ||
		commit.description != "")
}

// sendToMattermost : Envoi sur le Webhook de Mattermost
func sendToMattermost(repository string, commit commitInformation) {
	if !isCommitValid(commit) {
		err := errors.New("No Git repository found")
		toolbox.CheckError(err, 2)
	}

	// Récupération des paramètres
	// ---------------------------
	hookURL = config.MattermostHookURL
	hookPayload = config.MattermostHookPayload

	// Création du payload à transmettre
	// ---------------------------------
	payload := payload{
		Text: "",
	}
	payload.Text = "### Mise en production\n"
	payload.Text += "#### " + toolbox.Ucfirst(repository)

	if commit.version != "" {
		payload.Text += " - v" + commit.version
	}
	payload.Text += "\n"

	if commit.subject != "" {
		payload.Text += "| Sujet |" + toolbox.Ucfirst(commit.subject) + "|\n"
	}

	payload.Text += "|:---|:---|\n"

	if commit.author != "" {
		payload.Text += "| Auteur |" + toolbox.Ucfirst(commit.author) + "|\n"
	}

	if commit.developers != "" {
		payload.Text += "| Développeur(s) |" + commit.developers + "|\n"
	}

	if commit.testers != "" {
		payload.Text += "| Testeur(s) |" + commit.testers + "|\n"
	}

	if commit.description != "" {
		payload.Text += "#### Description :\n" + commit.description + "\n"
	}

	payloadJSONEncoded, err := json.Marshal(payload)
	toolbox.CheckError(err, 3)

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
