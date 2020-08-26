package server

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" //Mysql driver
)

//CreateConnection to database
func CreateConnection(connString string) (db *sql.DB)  {
	db, err := sql.Open("mysql", connString + "?parseTime=true")
	if err != nil {
		log.Print("Error: ", err)
	} else {
		log.Print("Connected to database")
	}
	return db
}