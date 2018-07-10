package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/fabienbellanger/goMattermost/models"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/labstack/echo"
)

// GetCommitsHandler : Liste des commits
func GetCommitsHandler(c echo.Context) error {
	// Limit du nombre de résultat
	// ---------------------------
	const limitMax = 50
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	toolbox.CheckError(err, 0)

	if limit > limitMax {
		limit = limitMax
	}

	// Tri
	// ---
	const defaultSort = "ASC"
	sort := strings.ToUpper(c.QueryParam("sort"))

	if sort != "ASC" && sort != "DESC" {
		sort = defaultSort
	}

	// Récupération des commits
	// ------------------------
	commits, err := models.GetCommitsList(limit, sort)
	toolbox.CheckError(err, 0)

	return c.JSON(http.StatusOK, commits)
}
