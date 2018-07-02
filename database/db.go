package database

import (
	"database/sql"

	"github.com/fabienbellanger/goMattermost/config"
	"github.com/fabienbellanger/goMattermost/toolbox"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var (
	// DB is the connection handle
	DB *sql.DB
)

// Open : Ouverture de la connexion
func Open() {
	db, err := sql.Open(config.DatabaseDriver,
		config.DatabaseUser+":"+config.DatabasePassword+"@/"+config.DatabaseName+"?parseTime=true")
	toolbox.CheckError(err, 0)

	DB = db
}

// prepareQuery : Préparation de la requête
func prepareQuery(query string) *sql.Stmt {
	statement, err := DB.Prepare(query)
	toolbox.CheckError(err, 0)

	return statement
}

// executeQuery : Exécute une requête de type INSERT, UPDATE ou DELETE
func executeQuery(query string, args ...interface{}) (sql.Result, error) {
	statement := prepareQuery(query)
	defer statement.Close()

	result, err := statement.Exec(args...)
	toolbox.CheckError(err, 0)

	return result, err
}

// Select : Exécution d'une requête
func Select(query string, args ...interface{}) (*sql.Rows, error) {
	statement := prepareQuery(query)
	defer statement.Close()

	rows, err := statement.Query(args...)
	toolbox.CheckError(err, 0)

	return rows, err
}

// Insert : Requête d'insertion
func Insert(query string, args ...interface{}) (int64, error) {
	result, err := executeQuery(query, args...)
	toolbox.CheckError(err, 0)

	id, err := result.LastInsertId()
	toolbox.CheckError(err, 0)

	return id, err
}

// Update : Requête de mise à jour
func Update(query string, args ...interface{}) (int64, error) {
	result, err := executeQuery(query, args...)
	toolbox.CheckError(err, 0)

	affect, err := result.RowsAffected()
	toolbox.CheckError(err, 0)

	return affect, err
}

// Delete : Requête de suppression
func Delete(query string, args ...interface{}) (int64, error) {
	result, err := executeQuery(query, args...)
	toolbox.CheckError(err, 0)

	affect, err := result.RowsAffected()
	toolbox.CheckError(err, 0)

	return affect, err
}
