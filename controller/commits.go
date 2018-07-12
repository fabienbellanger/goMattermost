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

	response := model.GetHTTPResponse(http.StatusOK, "Success", commits)

	return c.JSON(response.Code, response)
}

// GetCommitHandler : Récupération d'un commit
func GetCommitHandler(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	toolbox.CheckError(err, 0)

	commit := model.CommitJSON{}
	response := model.HTTPResponse{}

	if id != 0 {
		commit, err = model.GetCommit(id)

		if commit.ID == 0 {
			response.Code = http.StatusNotFound
			response.Message = "No commit found"
		} else {
			response.Code = http.StatusOK
			response.Message = "Success"
			response.Data = commit
		}
	} else {
		response.Code = http.StatusNotFound
		response.Message = "No commit found"
	}

	return c.JSON(response.Code, response)
}
