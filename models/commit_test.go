package models

import (
	"testing"

	"github.com/fabienbellanger/goMattermost/models"
)

// TestNewCommitDBFromCommitInformation : Test de la création d'une variable de type CommitDB
func TestNewCommitDBFromCommitInformation(t *testing.T) {
	repository := "repo"

	commit := models.CommitInformation{"author", "subject", "version", "developers", "testers", "description"}
	commitDB := newCommitDBFromCommitInformation(repository, commit)
	commitDBValid := models.CommitDB{}
}
