package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	// "github.com/golang-jwt/jwt"
	"github.com/r1chter1/d-backend/database"
	"github.com/r1chter1/d-backend/models"
)
var DBConn *sql.DB

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("Someone wants to log in")
}

func register(w http.ResponseWriter, r *http.Request) {
	var data map[string]string

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, v := range data {
		if len(v) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	user := models.User{
		Name:     data["name"],
		Phone:    data["phone"],
		Username: data["username"],
	}
	thisisapassword, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	user.Password = string(thisisapassword)
	log.Println(len(user.Password))
	log.Println(user)
	query := "INSERT INTO users(username, name, phone, password) VALUES($1, $2, $3, $4)"
	_, err = DBConn.Query(query, user.Username, user.Name, user.Phone, user.Password)
	CheckError(err)
	// log.Println(thisisapassword)
	log.Println(user)
}

func handleRequests() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	log.Println("------------------------------------------")
	log.Println("App is running on http://127.0.0.1:8000/")
	database.ConnectDb()
	handleRequests()
}
