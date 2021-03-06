package model

import (
	"database/sql"
	"errors"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/fabienbellanger/goMattermost/database"
	"github.com/fabienbellanger/goMattermost/toolbox"
)

// CommitDB type
type CommitDB struct {
	ID          uint64
	Project     string
	Version     sql.NullString
	Author      sql.NullString
	Subject     string
	Description sql.NullString
	Developers  sql.NullString
	Testers     sql.NullString
	CreatedAt   sql.RawBytes
}

// CommitJSON type
type CommitJSON struct {
	ID          uint64 `json:"id" xml:"id"`
	Project     string `json:"project" xml:"project"`
	Version     string `json:"version" xml:"version"`
	Author      string `json:"author" xml:"author"`
	Subject     string `json:"subject" xml:"subject"`
	Description string `json:"description" xml:"description"`
	Developers  string `json:"developers" xml:"developers"`
	Testers     string `json:"testers" xml:"testers"`
	CreatedAt   string `json:"createdAt" xml:"createdAt"`
}

// CommitInformation type
type CommitInformation struct {
	Author      string
	Subject     string
	Version     string
	Developers  string
	Testers     string
	Description string
}

// FormatRepositoryName : Formattage du nom du répository
func FormatRepositoryName(repository *string) {
	s := *repository
	lastIndex := strings.LastIndex(s, "/")

	repositoryLen := len(s)

	if lastIndex == repositoryLen-1 {
		*repository = s[0 : repositoryLen-1]
	}
}

// FormatGitCommit : Formattage du commit
func FormatGitCommit(gitLogOutput []byte, commit *CommitInformation) {
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
		regexMessage := regexp.MustCompile(`(?s)Version(?:\s)?: (.*)\nDev(?:\s)?: (.*)\n(?:Test|Tests)(?:\s)?: (.*)\nDescription(?:\s)?:\n(.*)`)
		regexResult := regexMessage.FindAllSubmatch([]byte(message), -1)

		if len(regexResult) == 0 {
			commit.Version = ""
			commit.Developers = ""
			commit.Testers = ""
			commit.Description = message
		} else {
			for _, matchMessage := range regexResult {
				if len(matchMessage) == 5 {
					commit.Version = string(matchMessage[1])
					commit.Developers = string(matchMessage[2])
					commit.Testers = string(matchMessage[3])
					commit.Description = string(matchMessage[4])
				} else {
					commit.Version = ""
					commit.Developers = ""
					commit.Testers = ""
					commit.Description = message
				}
			}
		}
	}
}

// RetrieveCommit : Récupération du dernier commit Git de la branche de production
func RetrieveCommit(path, branch string) []byte {
	// Si la branche de production n'est pas bien initialisée, on met master par défaut
	if branch == "" {
		branch = "master"
	}

	gitLogCommand := exec.Command("git",
		"log",
		"-1",
		"--pretty=format:<%an>%n<%s>%n<%b>",
		"--first-parent", branch,
	)
	gitLogCommand.Dir = path
	gitLogOutput, err := gitLogCommand.Output()
	toolbox.CheckError(err, 1)

	return gitLogOutput
}

// IsCommitValid : Les informations du commit sont-elles valides ?
func IsCommitValid(commit CommitInformation) bool {
	return (commit.Author != "" ||
		commit.Subject != "" ||
		commit.Version != "" ||
		commit.Developers != "" ||
		commit.Testers != "" ||
		commit.Description != "")
}

// newCommitDBFromCommitInformation : Création d'une variable de type CommitDB à partir d'une variable de type CommitInformation
func newCommitDBFromCommitInformation(repository string, commit CommitInformation) CommitDB {
	commitDB := CommitDB{}

	commitDB.Project = repository
	commitDB.Subject = commit.Subject

	commitDB.Version.String = commit.Version
	if len(commit.Version) == 0 {
		commitDB.Version.Valid = false
	} else {
		commitDB.Version.Valid = true
	}

	commitDB.Author.String = commit.Author
	if len(commit.Author) == 0 {
		commitDB.Author.Valid = false
	} else {
		commitDB.Author.Valid = true
	}

	commitDB.Description.String = commit.Description
	if len(commit.Description) == 0 {
		commitDB.Description.Valid = false
	} else {
		commitDB.Description.Valid = true
	}

	commitDB.Developers.String = commit.Developers
	if len(commit.Developers) == 0 {
		commitDB.Developers.Valid = false
	} else {
		commitDB.Developers.Valid = true
	}

	commitDB.Testers.String = commit.Testers
	if len(commit.Testers) == 0 {
		commitDB.Testers.Valid = false
	} else {
		commitDB.Testers.Valid = true
	}

	return commitDB
}

// AddCommit : Ajout d'un commit en base de données
func AddCommit(repository string, commit CommitInformation, commitDBChan chan<- CommitDB, errInsert chan<- error) {
	// Tests des données
	// -----------------
	if len(repository) == 0 || len(commit.Subject) == 0 {
		err := errors.New("Empty repository or empty subject")
		toolbox.CheckError(err, 1)
	}

	commitDB := newCommitDBFromCommitInformation(repository, commit)

	query := `
		INSERT INTO commit(project, version, author, subject, description, developers, testers, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`

	id, err := database.Insert(query,
		repository,
		commitDB.Version,
		commitDB.Author,
		commitDB.Subject,
		commitDB.Description,
		commitDB.Developers,
		commitDB.Testers)

	errInsert <- err

	if err == nil {
		commitDB.ID = uint64(id)
	}

	commitDBChan <- commitDB
}

// GetCommitsList : Liste des commits
func GetCommitsList(limit int, sort string) ([]CommitJSON, error) {
	query := `
		SELECT id, project, version, author, subject, description, developers, testers, created_at
		FROM commit
		ORDER BY created_at ` + sort + `
		LIMIT ?`
	rows, err := database.Select(query, limit)

	var commits = make([]CommitJSON, 0)
	var id uint64
	var project, subject string
	var version, author, description, developers, testers sql.NullString
	var createdAt sql.RawBytes

	for rows.Next() {
		err = rows.Scan(
			&id,
			&project,
			&version,
			&author,
			&subject,
			&description,
			&developers,
			&testers,
			&createdAt)

		datetime, _ := time.Parse(time.RFC3339, string(createdAt))

		commits = append(commits, CommitJSON{
			id,
			project,
			version.String,
			author.String,
			subject,
			description.String,
			developers.String,
			testers.String,
			datetime.Format("2006-01-02 15:04:05")})

		if err != nil {
			panic(err.Error())
		}
	}

	return commits, err
}

// GetCommit : Récupération d'un commit
func GetCommit(id int) (CommitJSON, error) {
	query := `
		SELECT id, project, version, author, subject, description, developers, testers, created_at
		FROM commit
		WHERE id = ?
		LIMIT 1`
	rows, err := database.Select(query, id)

	commit := CommitJSON{}
	var idDB uint64
	var project, subject string
	var version, author, description, developers, testers sql.NullString
	var createdAt sql.RawBytes

	for rows.Next() {
		err = rows.Scan(
			&idDB,
			&project,
			&version,
			&author,
			&subject,
			&description,
			&developers,
			&testers,
			&createdAt)

		datetime, _ := time.Parse(time.RFC3339, string(createdAt))

		commit.ID = idDB
		commit.Project = project
		commit.Version = version.String
		commit.Author = author.String
		commit.Subject = subject
		commit.Description = description.String
		commit.Developers = developers.String
		commit.Testers = testers.String
		commit.CreatedAt = datetime.Format("2006-01-02 15:04:05")

		if err != nil {
			panic(err.Error())
		}
	}

	return commit, err
}

// DeleteCommit : Suppression d'un commit
func DeleteCommit(id int) int64 {
	query := `DELETE FROM commit WHERE id = ?`
	nbDeleted, err := database.Delete(query, id)
	toolbox.CheckError(err, 0)

	return nbDeleted
}

// GetDailyCommitsForEmailing : Liste des commits à envoyer quotidiennement par email
func GetDailyCommitsForEmailing() []CommitJSON {
	/* dateNow := time.Now().Format("2006-01-02%")
	query := `
		SELECT project, version, author, subject, description, developers, testers, created_at
		FROM commit
		WHERE created_at LIKE ?`
	rows, err := database.Select(query, dateNow)*/
	query := `
		SELECT
			c1.project, c1.version, c1.author, c1.subject, c1.description,
			c1.developers, c1.testers, c1.created_at
		FROM commit c1
			INNER JOIN (
				SELECT MAX(DATE(created_at)) AS max
				FROM commit
			) c2 ON DATE(c1.created_at) = c2.max`
	rows, err := database.Select(query)
	toolbox.CheckError(err, 1)

	var commits = make([]CommitJSON, 0)
	var project, subject string
	var version, author, description, developers, testers sql.NullString
	var createdAt sql.RawBytes

	for rows.Next() {
		err = rows.Scan(
			&project,
			&version,
			&author,
			&subject,
			&description,
			&developers,
			&testers,
			&createdAt)

		datetime, _ := time.Parse(time.RFC3339, string(createdAt))

		commits = append(commits, CommitJSON{
			0,
			project,
			version.String,
			author.String,
			subject,
			description.String,
			developers.String,
			testers.String,
			datetime.Format("2006-01-02 15:04:05")})

		toolbox.CheckError(err, 1)
	}

	return commits
}
