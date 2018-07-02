package toolbox

import (
	"testing"
)

// TestCheckError : Test de la fonction CheckError
func TestCheckError(t *testing.T) {
	// TODO
}

// TestUcfirst : Test de la fonction Ucfirst
func TestUcfirst(t *testing.T) {
	s1 := "test chaîne Avec et sans majuscule"
	s1u := Ucfirst(s1)

	if s1u != "Test chaîne Avec et sans majuscule" {
		t.Errorf("Ucfirst - got: %q, want: %q.", s1u, "Test chaîne Avec et sans majuscule")
	}

	s2 := "Test chaîne Avec et sans majuscule"
	s2u := Ucfirst(s2)

	if s2u != "Test chaîne Avec et sans majuscule" {
		t.Errorf("Ucfirst - got: %q, want: %q.", s2u, "Test chaîne Avec et sans majuscule")
	}

	s3 := "été"
	s3u := Ucfirst(s3)

	if s3u != "Été" {
		t.Errorf("Ucfirst - got: %q, want: %q.", s3u, "Été")
	}
}
