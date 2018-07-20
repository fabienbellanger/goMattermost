package lib

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

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

// formattedCommit type for email
type formattedCommit struct {
	project    string
	version    string
	time       string
	developers []string
	testers    []string
	issues     []struct {
		action      string
		description string
	}
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

// SendCommitsByMail : Envoi les commits de la veille par email
func SendCommitsByMail() {
	// Récupération des commits de la veille
	// -------------------------------------
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	commits, err := model.GetDailyCommitsForEmailing(yesterday)
	toolbox.CheckError(err, 1)

	fmt.Println(commits)

	// Traitements des commits
	// -----------------------
	formattedCommits := formatCommits(commits)

	// Envoi du mail
	// -------------
	sendMail()
}

func formatCommits(commits model.CommitJSON) []formattedCommit {
	formattedCommits := make([]formattedCommit, 0)

	// TODO !!!

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
