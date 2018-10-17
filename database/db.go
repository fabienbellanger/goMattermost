package database

import (
	"database/sql"
	"os"
	"os/exec"
	"time"

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

// InitDatabase : Initialisation de la base de données
func InitDatabase() {
	// Requètes
	// --------
	queries := make([]string, 0)

	// User
	queries = append(queries, "DROP TABLE IF EXISTS user")
	queries = append(queries, `
		CREATE TABLE user (
			id int(10) unsigned NOT NULL AUTO_INCREMENT,
			username varchar(128) NOT NULL,
			password varchar(128) NOT NULL,
			lastname varchar(100) NOT NULL,
			firstname varchar(100) NOT NULL,
			created_at timestamp NULL DEFAULT NULL,
			deleted_at timestamp NULL DEFAULT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	// Commit
	queries = append(queries, "DROP TABLE IF EXISTS commit")
	queries = append(queries, `
		CREATE TABLE commit (
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			project varchar(50) NOT NULL,
			version varchar(11) DEFAULT '',
			author varchar(100) DEFAULT '',
			subject varchar(200) NOT NULL DEFAULT '',
			description text DEFAULT NULL,
			developers varchar(200) DEFAULT '',
			testers varchar(200) DEFAULT '',
			created_at datetime NOT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)

	transaction, err := DB.Begin()
	toolbox.CheckError(err, 1)

	defer func() {
		// Rollback the transaction after the function returns.
		// If the transaction was already commited, this will do nothing.
		_ = transaction.Rollback()
	}()

	for _, query := range queries {
		// Execute the query in the transaction.
		_, err := transaction.Exec(query)
		toolbox.CheckError(err, 1)
	}

	// Commit the transaction.
	err = transaction.Commit()
	toolbox.CheckError(err, 1)
}

// DumpDatabase : Dump de la base de données
func DumpDatabase() (string, int) {
	// Exécution de la commande
	// ------------------------
	dumpCommand := exec.Command("mysqldump",
		"-u"+config.DatabaseUser,
		"-p"+config.DatabasePassword,
		config.DatabaseName)
	dumpCommand.Dir = "."
	dumpOutput, err := dumpCommand.Output()
	toolbox.CheckError(err, 1)

	// Création du fichier
	// -------------------
	dumpFileName := "dump_" + time.Now().Format("2006-01-02_150405") + ".sql"
	file, err := os.Create(dumpFileName)
	toolbox.CheckError(err, 2)

	defer file.Close()

	size, err := file.Write(dumpOutput)
	toolbox.CheckError(err, 3)

	return dumpFileName, size
}
