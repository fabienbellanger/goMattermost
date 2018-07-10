package model

import (
	"testing"
)

// TestGetFullname : Test du nom complet d'un utilisateur
func TestGetFullname(t *testing.T) {
	user := UserDB{
		ID:        1,
		Lastname:  "Bellanger",
		Firstname: "Fabien",
	}

	fullname := user.GetFullname()
	userFullname := user.Firstname + " " + user.Lastname

	if fullname != userFullname {
		t.Errorf("getFullname - got: %s, want: %s.", fullname, userFullname)
	}
}
