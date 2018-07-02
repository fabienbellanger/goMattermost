package mattermost

import (
	"testing"
)

// TestIsCommitValid : Test de la fonction isCommitValid
func TestIsCommitValid(t *testing.T) {
	// Test avec un commit vide
	// ------------------------
	commitEmpty := commitInformation{}
	isCommitEmptyValid := isCommitValid(commitEmpty)

	if isCommitValid(commitEmpty) {
		t.Errorf("isCommitValid - got: %t, want: %t.", isCommitEmptyValid, false)
	}

	// Test avec un commit non vide
	// ----------------------------
	commit1 := commitInformation{}
	commit1.author = "Fabien Bellanger"
	commit1.version = "1.0.0"
	isCommit1Valid := isCommitValid(commit1)

	if !isCommit1Valid {
		t.Errorf("isCommitValid - got: %t, want: %t.", isCommit1Valid, true)
	}
}

// TestFormatGitCommit : Test du formattage du dernier commit sur master
func TestFormatGitCommit(t *testing.T) {
	// Commit complet
	gitlog1 := []byte("<Fabien>\n<Subject>\n<Message>")
	commit1 := commitInformation{}

	formatGitCommit(gitlog1, &commit1)

	if commit1.author == "" || commit1.subject == "" {
		t.Error("formatGitCommit - got: empty commit, want: commit with info.")
	}

	// Commit vide
	gitlog2 := []byte("<>\n<>\n<>")
	commit2 := commitInformation{}

	formatGitCommit(gitlog2, &commit2)

	if commit2.author != "" && commit2.subject != "" {
		t.Error("formatGitCommit - got: not empty commit, want: empty commit.")
	}
}
