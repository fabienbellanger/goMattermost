package mattermost

import (
	"testing"

	"github.com/fabienbellanger/goMattermost/models"
)

// TestIsCommitValid : Test de la fonction isCommitValid
func TestIsCommitValid(t *testing.T) {
	// Test avec un commit vide
	// ------------------------
	commitEmpty := models.CommitInformation{}
	isCommitEmptyValid := isCommitValid(commitEmpty)

	if isCommitValid(commitEmpty) {
		t.Errorf("isCommitValid - got: %t, want: %t.", isCommitEmptyValid, false)
	}

	// Test avec un commit non vide
	// ----------------------------
	commit1 := models.CommitInformation{}
	commit1.Author = "Fabien Bellanger"
	commit1.Version = "1.0.0"
	isCommit1Valid := isCommitValid(commit1)

	if !isCommit1Valid {
		t.Errorf("isCommitValid - got: %t, want: %t.", isCommit1Valid, true)
	}
}

// TestFormatGitCommit : Test du formattage du dernier commit sur master
func TestFormatGitCommit(t *testing.T) {
	// Commit complet
	gitlog1 := []byte("<Fabien>\n<Subject>\n<Message>")
	commit1 := models.CommitInformation{}

	formatGitCommit(gitlog1, &commit1)

	if commit1.Author == "" || commit1.Subject == "" {
		t.Error("formatGitCommit - got: empty commit, want: commit with info.")
	}

	// Commit vide
	gitlog2 := []byte("<>\n<>\n<>")
	commit2 := models.CommitInformation{}

	formatGitCommit(gitlog2, &commit2)

	if commit2.Author != "" && commit2.Subject != "" {
		t.Error("formatGitCommit - got: not empty commit, want: empty commit.")
	}
}

// TestFormatRepositoryName : Test du formattage du nom du répository
func TestFormatRepositoryName(t *testing.T) {
	rRef := "repo/"
	rGood := "repo"
	r := rRef
	formatRepositoryName(&r)

	if r != rGood {
		t.Errorf("formatRepositoryName - got: %s, want: %s.", r, rGood)
	}

	rRef = "repo"
	rGood = "repo"
	r = rRef
	formatRepositoryName(&r)

	if r != rGood {
		t.Errorf("formatRepositoryName - got: %s, want: %s.", r, rGood)
	}

	rRef = "/repo/sub/"
	rGood = "/repo/sub"
	r = rRef
	formatRepositoryName(&r)

	if r != rGood {
		t.Errorf("formatRepositoryName - got: %s, want: %s.", r, rGood)
	}
}
