package model

import (
	"testing"
)

// TestIsCommitValid : Test de la fonction IsCommitValid
func TestIsCommitValid(t *testing.T) {
	// Test avec un commit vide
	// ------------------------
	commitEmpty := CommitInformation{}
	isCommitEmptyValid := IsCommitValid(commitEmpty)

	if IsCommitValid(commitEmpty) {
		t.Errorf("IsCommitValid - got: %t, want: %t.", isCommitEmptyValid, false)
	}

	// Test avec un commit non vide
	// ----------------------------
	commit1 := CommitInformation{}
	commit1.Author = "Fabien Bellanger"
	commit1.Version = "1.0.0"
	isCommit1Valid := IsCommitValid(commit1)

	if !isCommit1Valid {
		t.Errorf("IsCommitValid - got: %t, want: %t.", isCommit1Valid, true)
	}
}

// TestFormatGitCommit : Test du formattage du dernier commit sur master
func TestFormatGitCommit(t *testing.T) {
	// Commit complet
	gitlog1 := []byte("<Fabien>\n<Subject>\n<Message>")
	commit1 := CommitInformation{}

	FormatGitCommit(gitlog1, &commit1)

	if commit1.Author == "" || commit1.Subject == "" {
		t.Error("FormatGitCommit - got: empty commit, want: commit with info.")
	}

	// Commit vide
	gitlog2 := []byte("<>\n<>\n<>")
	commit2 := CommitInformation{}

	FormatGitCommit(gitlog2, &commit2)

	if commit2.Author != "" && commit2.Subject != "" {
		t.Error("FormatGitCommit - got: not empty commit, want: empty commit.")
	}
}

// TestFormatRepositoryName : Test du formattage du nom du repository
func TestFormatRepositoryName(t *testing.T) {
	rRef := "repo/"
	rGood := "repo"
	r := rRef
	FormatRepositoryName(&r)

	if r != rGood {
		t.Errorf("FormatRepositoryName - got: %s, want: %s.", r, rGood)
	}

	rRef = "repo"
	rGood = "repo"
	r = rRef
	FormatRepositoryName(&r)

	if r != rGood {
		t.Errorf("FormatRepositoryName - got: %s, want: %s.", r, rGood)
	}

	rRef = "/repo/sub/"
	rGood = "/repo/sub"
	r = rRef
	FormatRepositoryName(&r)

	if r != rGood {
		t.Errorf("FormatRepositoryName - got: %s, want: %s.", r, rGood)
	}
}

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
