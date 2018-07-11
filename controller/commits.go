package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/fabienbellanger/goMattermost/model"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/labstack/echo"
)

// GetCommitsHandler : Liste des commits
func GetCommitsHandler(c echo.Context) error {
	// Limit du nombre de résultat
	// ---------------------------
	const limitMax = 50
	limitParam := c.QueryParam("limit")
	limit := limitMax

	if limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		toolbox.CheckError(err, 0)
	}

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
	commits, err := model.GetCommitsList(limit, sort)
	toolbox.CheckError(err, 0)

	return c.JSON(http.StatusOK, commits)
}

// GetCommitHandler : Récupération d'un commit
func GetCommitHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	toolbox.CheckError(err, 0)

	if id != 0 {

	}

	var err1 error

	return err1
}
