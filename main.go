package main

import (
	"fmt"
	"log"
	"net/http"
	// "github.com/golang-jwt/jwt"
	"github.com/r1chter1/d-backend/database"
) 

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Someone wants to log in")
}

func register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Someone wants to register")
}

func handleRequests() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	fmt.Println("App is running on http://127.0.0.1:8000/")
	connectDb()
	handleRequests()
}
