package database

import (
	"database/sql"
	"fmt"
	"log"
)

var DBConn *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "bekbull"
	password = "ident"
	dbname   = "tumar"
)

func init() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	// open database
	DBConn, err = sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	log.Println("Connected!")
}
