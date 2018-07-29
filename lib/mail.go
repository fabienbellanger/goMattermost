package lib

import (
	"fmt"
	"net/smtp"
	"regexp"
	"strings"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/model"
	"github.com/fabienbellanger/goMattermost/toolbox"
)

// Mail type
type Mail struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Body    string
}

// issue type
type issue struct {
	action      string
	description string
}

// formattedCommit type for email
type formattedCommit struct {
	project    string
	version    string
	time       string
	developers []string
	testers    []string
	issues     []issue
}

// serverName : Nom du serveur
func serverName() (s string) {
	if len(config.SMTPPort) > 0 {
		s = config.SMTPHost + ":" + config.SMTPPort
	} else {
		s = config.SMTPHost
	}

	return
}

// buildMessage : Construction du body
func (mail *Mail) buildMessage() (header string) {
	header = ""
	header += fmt.Sprintf("From: %s\r\n", mail.From)

	if len(mail.To) > 0 {
		header += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	}

	if len(mail.Cc) > 0 {
		header += fmt.Sprintf("Cc: %s\r\n", strings.Join(mail.Cc, ";"))
	}

	header += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	header += fmt.Sprintf("MIME-version: %s\r\n", "1.0")
	header += fmt.Sprintf("Content-Type: %s\r\n", "text/html")
	header += fmt.Sprintf("charset: %s\r\n", "UTF-8")
	header += "\r\n" + mail.Body

	return
}

// SendCommitsByMail : Envoi les commits du dernier jour par email
func SendCommitsByMail() {
	// Récupération des commits du dernier jour
	// ----------------------------------------
	commits := model.GetDailyCommitsForEmailing()
	fmt.Println(commits)

	// Traitements des commits
	// -----------------------
	formattedCommits := formatCommits(commits)
	fmt.Println(formattedCommits)

	// Envoi du mail
	// -------------
	// sendMail()
}

// formatCommits : Formattage des commits
func formatCommits(commits []model.CommitJSON) []formattedCommit {
	formattedCommits := make([]formattedCommit, 0)
	regexDescription := regexp.MustCompile(`- (?:\[(fix|add|improvement)\] )?(.*)`)
	developersTestersDelimiter := " & "

	var formattedCommit formattedCommit
	var issue issue

	for _, commit := range commits {
		formattedCommit.project = commit.Project
		formattedCommit.version = commit.Version
		formattedCommit.time = commit.CreatedAt[11:16]
		formattedCommit.developers = strings.Split(commit.Developers, developersTestersDelimiter)
		formattedCommit.testers = strings.Split(commit.Testers, developersTestersDelimiter)

		// Description
		matches := regexDescription.FindAllSubmatch([]byte(commits[0].Description), -1)
		for _, match := range matches {
			if len(match) == 3 {
				issue.action = string(match[1])
				issue.description = string(match[2])

				formattedCommit.issues = append(formattedCommit.issues, issue)
			}
		}

		formattedCommits = append(formattedCommits, formattedCommit)
	}

	return formattedCommits
}

// sendMail : Envoi du mail
func sendMail() {
	mail := Mail{}
	mail.From = "toto@hjdhs.fr"
	mail.To = []string{"def@yahoo.com", "xyz@outlook.com"}
	mail.Cc = []string{"mnp@gmail.com"}
	mail.Bcc = []string{"a69@outlook.com"}
	mail.Subject = "Test envoi mails go"
	mail.Body = "This is the <b>email</b> body."

	messageBody := mail.buildMessage()

	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)
	err := smtp.SendMail(serverName(), auth, mail.From, mail.To, []byte(messageBody))
	toolbox.CheckError(err, 1)
}
