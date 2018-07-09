package controller

import (
	"net/http"
	"strconv"

	"github.com/fabienbellanger/goMattermost/models"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/labstack/echo"
)

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

	return c.JSON(http.StatusOK, commits)
}
