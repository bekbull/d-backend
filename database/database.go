package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "bekbull"
	password = "ident"
	dbname   = "tumar"
)

func ConnectDb() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	DBConn, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	// close database
	defer DBConn.Close()

	// check db
	err = DBConn.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Connected!")
}
