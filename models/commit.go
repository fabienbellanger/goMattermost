package models

import (
	"database/sql"
	"fmt"

	"github.com/fabienbellanger/goMattermost/database"
)

// CommitDB type
type CommitDB struct {
	ID          uint64
	Version     string
	Author      string
	Subject     string
	Description string
	Developers  string
	Testers     string
	CreatedAt   sql.RawBytes
}

// CommitInformation structure
type CommitInformation struct {
	Author      string
	Subject     string
	Version     string
	Developers  string
	Testers     string
	Description string
}

// AddCommit : Ajout d'un commit en base de données
func AddCommit(commit CommitInformation) {
	query := `
		INSERT INTO commit(version, author, subject, description, developers, testers, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())`

	fmt.Println(query)
	id, err := database.Insert(query,
		commit.Version,
		commit.Author,
		commit.Subject,
		commit.Description,
		commit.Developers,
		commit.Testers)

	fmt.Println(id, err)
}
