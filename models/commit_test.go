package models

import (
	"testing"
)

// TestNewCommitDBFromCommitInformation : Test de la création d'une variable de type CommitDB
func TestNewCommitDBFromCommitInformation(t *testing.T) {
	repository := "repo"
	commit := CommitInformation{"author", "subject", "version", "developers", "testers", "description"}
	commitDB := newCommitDBFromCommitInformation(repository, commit)

	commitDBValid := CommitDB{}
	commitDBValid.ID = 0
	commitDBValid.Project = "repo"
	commitDBValid.Version.String = "version"
	commitDBValid.Version.Valid = true
	commitDBValid.Author.String = "author"
	commitDBValid.Author.Valid = true
	commitDBValid.Subject = "subject"
	commitDBValid.Description.String = "description"
	commitDBValid.Description.Valid = true
	commitDBValid.Developers.String = "developers"
	commitDBValid.Developers.Valid = true
	commitDBValid.Testers.String = "testers"
	commitDBValid.Testers.Valid = true

	if commitDB.ID != commitDBValid.ID || commitDB.Project != commitDBValid.Project ||
		commitDB.Version != commitDBValid.Version || commitDB.Author != commitDBValid.Author ||
		commitDB.Subject != commitDBValid.Subject || commitDB.Description != commitDBValid.Description ||
		commitDB.Developers != commitDBValid.Developers || commitDB.Testers != commitDBValid.Testers {
		t.Errorf("newCommitDBFromCommitInformation - got %+v: , want: %+v.", commitDB, commitDBValid)
	}

	commit = CommitInformation{"", "subject", "", "developers", "testers", "description"}
	commitDB = newCommitDBFromCommitInformation(repository, commit)

	commitDBValid = CommitDB{}
	commitDBValid.ID = 0
	commitDBValid.Project = "repo"
	commitDBValid.Version.String = ""
	commitDBValid.Version.Valid = false
	commitDBValid.Author.String = ""
	commitDBValid.Author.Valid = false
	commitDBValid.Subject = "subject"
	commitDBValid.Description.String = "description"
	commitDBValid.Description.Valid = true
	commitDBValid.Developers.String = "developers"
	commitDBValid.Developers.Valid = true
	commitDBValid.Testers.String = "testers"
	commitDBValid.Testers.Valid = true

	if commitDB.ID != commitDBValid.ID || commitDB.Project != commitDBValid.Project ||
		commitDB.Version != commitDBValid.Version || commitDB.Author != commitDBValid.Author ||
		commitDB.Subject != commitDBValid.Subject || commitDB.Description != commitDBValid.Description ||
		commitDB.Developers != commitDBValid.Developers || commitDB.Testers != commitDBValid.Testers {
		t.Errorf("newCommitDBFromCommitInformation - got %+v: , want: %+v.", commitDB, commitDBValid)
	}
}
