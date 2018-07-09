package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fabienbellanger/goMattermost/models"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/labstack/echo"
)

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

// GetCommitsHandler : Liste des commits
func GetCommitsHandler(c echo.Context) error {
	// Limit du nombre de résultat
	const limitMax = 50
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	toolbox.CheckError(err, 1)

	if limit > limitMax {
		limit = limitMax
	}

	commits, err := models.GetCommitsList(limit)
	toolbox.CheckError(err, 1)

	commitsJSON := make([]CommitJSON, len(commits))

	// Conversion CommitDB en commitJSON
	// ---------------------------------
	for i, commit := range commits {
		datetime, _ := time.Parse(time.RFC3339, string(commit.CreatedAt))

		commitsJSON[i].ID = commit.ID
		commitsJSON[i].Project = commit.Project
		commitsJSON[i].Version = commit.Version.String
		commitsJSON[i].Author = commit.Author.String
		commitsJSON[i].Subject = commit.Subject
		commitsJSON[i].Description = commit.Description.String
		commitsJSON[i].Developers = commit.Developers.String
		commitsJSON[i].Testers = commit.Testers.String
		commitsJSON[i].CreatedAt = datetime.Format("2006-01-02 15:04:05")
	}

	return c.JSON(http.StatusOK, commitsJSON)
}
