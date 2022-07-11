// Info and funcs to connect to data sources (sqlite db, mock merit api, potentially others)

package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Harcoded for now - should be externalized to properties file, etc. 
var MeritMockApiBaseURL = "http://localhost:3000"
var dbPath = "C:/sqlite_dbs/beiq.db"

// sqlite db connection, available to all models
var DbConn *sql.DB

func ConnectDb () (error) {
	db, err := sql.Open("sqlite3", dbPath)
	//db, err := sql.Open("sqlite3", "C:/sqlite_dbs/beiq.db")
	if err != nil {
		return err
	}
	DbConn = db
	return nil
}

// Generic function to check errors with various db interaction, etc., in order to eliminate some redundant code
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}