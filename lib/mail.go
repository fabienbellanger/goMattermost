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

// Issue type
type Issue struct {
	Action      string
	Description string
}

// formattedCommit type for email
type formattedCommit struct {
	Project    string
	Version    string
	Time       string
	Developers []string
	Testers    []string
	Issues     []Issue
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

// mailTemplate type pour la template HTML
type mailTemplate struct {
	Title   string
	Commits []formattedCommit
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
	// formattedCommits := formatCommits(commits)
	// fmt.Println(formattedCommits)

	// Construction des données pour envoi à la template
	// -------------------------------------------------
	projects := constructData(commits)
	fmt.Println(projects)

	// Affiche les commits groupés par projet
	// --------------------------------------
	// mailbody := constructTemplate(formattedCommits)
	// fmt.Println(mailbody)

	// Envoi du mail
	// -------------
	// sendMail(mailbody)
}

// constructData : Construction des données pour envoi à la template
func constructData(commits []model.CommitJSON) []Project {
	projects := make([]Project, 0)
	regexDescription := regexp.MustCompile(`- (?:\[(fix|add|improvement|other)\] )?(.*)`)
	developersTestersDelimiter := " & "

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
			indexProject++

			projects = append(projects, Project{commit.Project, make([]Release, 0)})
		} else {
			indexProject = i
		}

		// 2. Regroupement par release
		// ---------------------------
		releasesNumber = len(projects[indexProject].Releases)
		indexRelease = releasesNumber - 1
		i = 0

		for i < releasesNumber && projects[indexProject].Releases[i].Version != commit.Version {
			i++
		}

		// La release est-elle déjà présente dans le tableau ?
		if i == releasesNumber {
			// Release non trouvée
			indexRelease++

			projects[indexProject].Releases = append(projects[indexProject].Releases, Release{
				commit.Version,
				"Time",
				strings.Split(commit.Developers, developersTestersDelimiter),
				strings.Split(commit.Testers, developersTestersDelimiter),
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
		} else {
			return false
		}
	})

	// Tri par release pour chaque projet
	// ----------------------------------
	for index := range projects {
		sort.Slice(projects[index].Releases, func(i, j int) bool {
			if projects[index].Releases[i].Version <= projects[index].Releases[j].Version {
				return true
			} else {
				return false
			}
		})
	}

	return projects
}

// formatCommits : Formattage des commits
func formatCommits(commits []model.CommitJSON) []formattedCommit {
	formattedCommits := make([]formattedCommit, 0)
	regexDescription := regexp.MustCompile(`- (?:\[(fix|add|improvement|other)\] )?(.*)`)
	developersTestersDelimiter := " & "

	var formattedCommit formattedCommit
	var issue Issue

	for _, commit := range commits {
		formattedCommit.Project = commit.Project
		formattedCommit.Version = commit.Version
		formattedCommit.Time = commit.CreatedAt[11:16]
		formattedCommit.Developers = strings.Split(commit.Developers, developersTestersDelimiter)
		formattedCommit.Testers = strings.Split(commit.Testers, developersTestersDelimiter)

		// Description
		matches := regexDescription.FindAllSubmatch([]byte(commit.Description), -1)
		for _, match := range matches {
			if len(match) == 3 {
				issue.Action = string(match[1])
				issue.Description = string(match[2])

				formattedCommit.Issues = append(formattedCommit.Issues, issue)
			}
		}

		formattedCommits = append(formattedCommits, formattedCommit)
	}

	// Tri du tableau par projet puis par version
	// ------------------------------------------
	sort.Slice(formattedCommits, func(i, j int) bool {
		if formattedCommits[i].Project < formattedCommits[j].Project {
			return true
		} else if formattedCommits[i].Project > formattedCommits[j].Project {
			return false
		} else {
			return formattedCommits[i].Version < formattedCommits[j].Version
		}
	})

	// Construction du tableau final
	// -----------------------------

	return formattedCommits
}

// constructTemplate : Construction de la template pour l'envoi du mail
// TODO: https://github.com/mlabouardy/go-html-email
func constructTemplate(commits []formattedCommit) string {
	templateData := mailTemplate{Title: "Titre de ma page", Commits: commits}

	t := template.New("mail")
	t = template.Must(t.ParseFiles("./templates/mail.html"))

	buffer := new(bytes.Buffer)
	// err := t.ExecuteTemplate(buffer, "mail", m)
	err := t.Execute(buffer, templateData)

	if err != nil {
		fmt.Println(err)
	}

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

	fmt.Print(" -> Mail send: \t")
	color.Green("Success\n\n")
}
