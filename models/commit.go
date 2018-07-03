package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/fabienbellanger/goMattermost/database"
	"github.com/fabienbellanger/goMattermost/toolbox"
	"github.com/fatih/color"
)

// CommitDB type
type CommitDB struct {
	ID          uint64
	Project     string
	Version     sql.NullString
	Author      sql.NullString
	Subject     string
	Description sql.NullString
	Developers  sql.NullString
	Testers     sql.NullString
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

// newCommitDBFromCommitInformation : Création d'une variable de type CommitDB
// à partir d'une variable de type CommitInformation
func newCommitDBFromCommitInformation(repository string, commit CommitInformation) CommitDB {
	commitDB := CommitDB{}

	commitDB.Project = repository
	commitDB.Subject = commit.Subject

	commitDB.Version.String = commit.Version
	if len(commit.Version) == 0 {
		commitDB.Version.Valid = false
	} else {
		commitDB.Version.Valid = true
	}

	commitDB.Author.String = commit.Author
	if len(commit.Author) == 0 {
		commitDB.Author.Valid = false
	} else {
		commitDB.Author.Valid = true
	}

	commitDB.Description.String = commit.Description
	if len(commit.Description) == 0 {
		commitDB.Description.Valid = false
	} else {
		commitDB.Description.Valid = true
	}

	commitDB.Developers.String = commit.Developers
	if len(commit.Developers) == 0 {
		commitDB.Developers.Valid = false
	} else {
		commitDB.Developers.Valid = true
	}

	commitDB.Testers.String = commit.Testers
	if len(commit.Testers) == 0 {
		commitDB.Testers.Valid = false
	} else {
		commitDB.Testers.Valid = true
	}

	return commitDB
}

// AddCommit : Ajout d'un commit en base de données
func AddCommit(repository string, commit CommitInformation) {
	// Tests des données
	// -----------------
	if len(repository) == 0 || len(commit.Subject) == 0 {
		err := errors.New("Empty repository or empty subject")
		toolbox.CheckError(err, 1)
	}

	commitDB := newCommitDBFromCommitInformation(repository, commit)

	query := `
		INSERT INTO commit(project, version, author, subject, description, developers, testers, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`

	id, err := database.Insert(query,
		repository,
		commitDB.Version,
		commitDB.Author,
		commitDB.Subject,
		commitDB.Description,
		commitDB.Developers,
		commitDB.Testers)

	if err == nil {
		commitDB.ID = uint64(id)

		fmt.Print("Insertion dans la table commit avec l'ID : ")
		color.Green(strconv.FormatInt(id, 10) + "\n\n")
	}
}
