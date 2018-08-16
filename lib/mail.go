package lib

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"regexp"
	"sort"
	"strings"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/model"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
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

// mailTemplate type pour la template HTML
type mailTemplate struct {
	Title   string
	Commits []string
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

	// Traitements des commits
	// -----------------------
	formattedCommits := formatCommits(commits)

	// Affiche les commits groupés par projet
	// --------------------------------------
	mailbody := printCommits(formattedCommits)

	constructTemplate()

	// Envoi du mail
	// -------------
	sendMail(mailbody)
}

// formatCommits : Formattage des commits
func formatCommits(commits []model.CommitJSON) []formattedCommit {
	formattedCommits := make([]formattedCommit, 0)
	regexDescription := regexp.MustCompile(`- (?:\[(fix|add|improvement|other)\] )?(.*)`)
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
		matches := regexDescription.FindAllSubmatch([]byte(commit.Description), -1)
		for _, match := range matches {
			if len(match) == 3 {
				issue.action = string(match[1])
				issue.description = string(match[2])

				formattedCommit.issues = append(formattedCommit.issues, issue)
			}
		}

		formattedCommits = append(formattedCommits, formattedCommit)
	}

	// Tri du tableau par projet puis par version
	// ------------------------------------------
	sort.Slice(formattedCommits, func(i, j int) bool {
		if formattedCommits[i].project < formattedCommits[j].project {
			return true
		} else if formattedCommits[i].project > formattedCommits[j].project {
			return false
		} else {
			return formattedCommits[i].version < formattedCommits[j].version
		}
	})

	// Construction du tableau final
	// -----------------------------

	return formattedCommits
}

// printCommits : Affichage des commits
func printCommits(commits []formattedCommit) string {
	var project string
	var color string
	var message string

	str := ""
	for index, commit := range commits {
		if commit.project != project {
			project = commit.project

			if index > 0 {
				str += "</ul>"
			}
			str += "<p style=\"font-weight: bold\">" + toolbox.Ucfirst(project) + "</p>"
			str += "<ul>"
		}

		str += "<li>"
		str += "[" + commit.version + "] [" + commit.time + "] "

		// Message
		message = ""
		for _, issue := range commit.issues {
			if issue.description != "" {
				if issue.action == "fix" {
					color = "red"
				} else if issue.action == "improvement" {
					color = "orange"
				} else if issue.action == "add" {
					color = "green"
				} else {
					color = "black"
				}

				message += "<li style=\"color: " + color + "\">"
				if issue.action != "" {
					message += "[" + issue.action + "] "
				}
				message += issue.description
				message += "</li>"
			}
		}

		if message != "" {
			str += "<ul>" + message + "</ul>"
		}

		str += "</li>"
	}
	str += "<ul>"

	return str
}

// constructTemplate : Construction de la template pour l'envoi du mail
// TODO: https://github.com/mlabouardy/go-html-email
func constructTemplate() {
	// Création d'une page
	c := []string{"coucou", "Toto"}
	m := mailTemplate{Title: "Titre de ma page", Commits: c}

	t := template.New("mail")
	t = template.Must(t.ParseFiles("./templates/mail.html"))

	buffer := new(bytes.Buffer)
	// err := t.ExecuteTemplate(buffer, "mail", m)
	err := t.Execute(buffer, m)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(buffer.String())
}

// sendMail : Envoi du mail
func sendMail(body string) {
	mail := Mail{}
	mail.From = "toto@hjdhs.fr"
	mail.To = []string{"def@yahoo.com", "xyz@outlook.com"}
	mail.Cc = []string{"mnp@gmail.com"}
	mail.Bcc = []string{"a69@outlook.com"}
	mail.Subject = "Test envoi mails go"
	mail.Body = body

	messageBody := mail.buildMessage()

	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)
	err := smtp.SendMail(serverName(), auth, mail.From, mail.To, []byte(messageBody))
	toolbox.CheckError(err, 1)

	fmt.Print(" -> Mail send: \t")
	color.Green("Success\n\n")
}
