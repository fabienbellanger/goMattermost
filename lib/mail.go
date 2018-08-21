package lib

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"regexp"
	"sort"
	"strings"
	"time"

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

// Project type
type Project struct {
	Name     string
	Releases []Release
}

// Release type
type Release struct {
	Version    string
	Time       string
	Developers []string
	Testers    []string
	Issues     []Issue
}

// Issue type
type Issue struct {
	Action      string
	Label       string
	Description string
}

// mailTemplate type pour la template HTML
type mailTemplate struct {
	Date     string
	Projects []Project
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

	if len(commits) > 0 {
		// Construction des données pour envoi à la template
		// -------------------------------------------------
		projects := constructData(commits)

		// Affiche les commits groupés par projet
		// --------------------------------------
		mailbody := constructTemplate(projects)

		// Envoi du mail
		// -------------
		sendMail(mailbody)
	}
}

// constructData : Construction des données pour envoi à la template
func constructData(commits []model.CommitJSON) []Project {
	projects := make([]Project, 0)
	regexDescription := regexp.MustCompile(`- (?:\[(fix|add|improvement|other)\] )?(.*)`)

	indexProject := -1
	indexRelease := -1
	projectsNumber := 0
	releasesNumber := 0
	var i int

	for _, commit := range commits {
		// 1. Regroupement par projet
		// --------------------------
		projectsNumber = len(projects)
		i = 0

		for i < projectsNumber && projects[i].Name != commit.Project {
			i++
		}

		// Le projet est-il déjà présent dans le tableau ?
		if i == projectsNumber {
			// Projet non trouvé
			indexProject = projectsNumber

			projects = append(projects, Project{commit.Project, make([]Release, 0)})
		} else {
			indexProject = i
		}

		// 2. Regroupement par release
		// ---------------------------
		releasesNumber = len(projects[indexProject].Releases)
		i = 0

		for i < releasesNumber && projects[indexProject].Releases[i].Version != commit.Version {
			i++
		}

		// La release est-elle déjà présente dans le tableau ?
		if i == releasesNumber {
			// Release non trouvée
			indexRelease = releasesNumber

			projects[indexProject].Releases = append(projects[indexProject].Releases, Release{
				commit.Version,
				commit.CreatedAt[11:16],
				getDevelopersTesters(commit.Developers),
				getDevelopersTesters(commit.Testers),
				make([]Issue, 0),
			})
		} else {
			indexRelease = i
		}

		// 3. Regroupement par issue
		// -------------------------
		// Traitement de la description
		matches := regexDescription.FindAllSubmatch([]byte(commit.Description), -1)
		for _, match := range matches {
			if len(match) == 3 {

				projects[indexProject].Releases[indexRelease].Issues = append(projects[indexProject].Releases[indexRelease].Issues, Issue{
					string(match[1]),
					getIssueLabel(string(match[1])),
					string(match[2]),
				})
			}
		}
	}

	// Tri du tableau par projet
	// -------------------------
	sort.Slice(projects, func(i, j int) bool {
		if projects[i].Name <= projects[j].Name {
			return true
		}

		return false
	})

	// Tri par release pour chaque projet
	// ----------------------------------
	for index := range projects {
		projects[index].Name = toolbox.Ucfirst(projects[index].Name)
		sort.Slice(projects[index].Releases, func(i, j int) bool {
			if projects[index].Releases[i].Version <= projects[index].Releases[j].Version {
				return true
			}

			return false
		})
	}

	return projects
}

// getIssueLabel : Retourne le label d'une issue
func getIssueLabel(action string) string {
	switch action {
	case "fix":
		return "Correction"
	case "improvement":
		return "Amélioration"
	case "add":
		return "Nouveauté"
	default:
		return "Autre"
	}
}

// getDevelopersTesters
func getDevelopersTesters(list string) []string {
	delimiter := " & "

	if list != "" {
		return strings.Split(list, delimiter)
	}

	return nil
}

// constructTemplate : Construction de la template pour l'envoi du mail
func constructTemplate(projects []Project) string {
	templateData := mailTemplate{Date: time.Now().Format("02/01/2006"), Projects: projects}

	t := template.New("mail")
	t = template.Must(t.ParseFiles("./templates/header.tmpl", "./templates/footer.tmpl", "./templates/mail.tmpl"))

	buffer := new(bytes.Buffer)
	err := t.Execute(buffer, templateData)
	toolbox.CheckError(err, 0)

	return buffer.String()
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

	fmt.Print(" -> Mail send:\t")
	color.Green("Success\n\n")
}
