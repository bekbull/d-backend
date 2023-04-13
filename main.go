package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/justinas/alice"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"

	"github.com/r1chter1/d-backend/database"
	"github.com/r1chter1/d-backend/models"
)

type UserIdKey struct{}

func login(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	var user models.User
	var query string
	var id int
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query = "SELECT username, password, id FROM users WHERE username=$1"
	err = database.DBConn.QueryRow(query, data["username"]).Scan(&user.Username, &user.Password, &id)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
	if err != nil {
		resp := make(map[string]string)
		resp["message"] = "Wrong Ppssword"
		w.WriteHeader(http.StatusExpectationFailed)
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := make(map[string]string)
	token, err := CreateToken(uint64(id))
	if err != nil {
		panic(err)
	}
	resp["token"] = token
	json.NewEncoder(w).Encode(resp)

	log.Println("Someone wants to log in")
}

func register(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	var query string
	var user models.User
	resp := make(map[string]string)
	w.Header().Set("Content-Type", "application/json")
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
	query = "SELECT username, phone FROM users WHERE username=$1 OR phone =$2"
	err = database.DBConn.QueryRow(query, data["username"], data["phone"]).Scan(&user.Username, &user.Phone)
	if err == sql.ErrNoRows {
		user = models.User{
			Name:     data["name"],
			Phone:    data["phone"],
			Username: data["username"],
		}
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
		user.Password = string(hashedPassword)

		query = "INSERT INTO users(username, name, phone, password) VALUES($1, $2, $3, $4)"
		_, err = database.DBConn.Exec(query, user.Username, user.Name, user.Phone, user.Password)
		CheckError(err)
		resp["message"] = "User created"
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if len(user.Username) > 0 {
		log.Println(user.Username)
		resp["message"] = "User already exists"
		w.WriteHeader(http.StatusExpectationFailed)
		json.NewEncoder(w).Encode(resp)
		return
	}

	log.Println(user)
}

func cart(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserIdKey{}).(string)
	response := make(map[string]string)

	log.Println(user)
	if r.Method == http.MethodGet {
		query := "SELECT * FROM cart WHERE username=$1"
		rows, err := database.DBConn.Query(query, user)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		// An items slice to hold data from returned rows.
		var items []models.Cart

		// Loop through rows, using Scan to assign column data to struct fields.
		for rows.Next() {
			var item models.Cart
			if err := rows.Scan(&item.IdNumber, &item.Username); err != nil {
				panic(err)
			}
			items = append(items, item)
		}
		if err = rows.Err(); err != nil {
			response["message"] = "Missing auth token"
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}
		json.NewEncoder(w).Encode(items)
	} else if r.Method == http.MethodPost {
		var data map[string]string

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		query := "INSERT INTO cart (code, username) VALUES($1, $2)"
		_, err = database.DBConn.Exec(query, data["code"], user)
		CheckError(err)
		response["message"] = "Item has been added"
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	} else if r.Method == http.MethodDelete {
		var data map[string]string

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		query := "DELETE FROM cart WHERE username=$1 code=$2"
		_, err = database.DBConn.Exec(query, user, data["code"])
		CheckError(err)
		response["message"] = "Item has been deleted"
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func order(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserIdKey{}).(string)
	response := make(map[string]string)

	log.Println(user)
	if r.Method == http.MethodGet {
		query := "SELECT * FROM order WHERE username=$1"
		rows, err := database.DBConn.Query(query, user)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		// An items slice to hold data from returned rows.
		var items []models.Cart

		// Loop through rows, using Scan to assign column data to struct fields.
		for rows.Next() {
			var item models.Cart
			if err := rows.Scan(&item.IdNumber, &item.Username); err != nil {
				panic(err)
			}
			items = append(items, item)
		}
		if err = rows.Err(); err != nil {
			response["message"] = "Missing auth token"
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}
		json.NewEncoder(w).Encode(items)
	} else if r.Method == http.MethodPost {
		var data map[string]string

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		query := "INSERT INTO order (code, username) VALUES($1, $2)"
		_, err = database.DBConn.Exec(query, data["code"], user)
		CheckError(err)
		response["message"] = "Item has been added"
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	} else if r.Method == http.MethodDelete {
		var data map[string]string

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		query := "DELETE FROM cart WHERE username=$1 code=$2"
		_, err = database.DBConn.Exec(query, user, data["code"])
		CheckError(err)
		response["message"] = "Item has been deleted"
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func handleRequests() {
	commonHandlers := alice.New(JwtAuthentication)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.Handle("/cart", commonHandlers.ThenFunc(cart))
	http.Handle("/order", commonHandlers.ThenFunc(order))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateToken(userid uint64) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["UserId"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 60 * 24 * 7).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	log.Println(os.Getenv("ACCESS_SECRET"))
	if err != nil {
		return "", err
	}
	return token, nil
}

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		response := make(map[string]string)
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response["message"] = "Missing auth token"
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response["message"] = "Something went wrong"
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			response["message"] = "Malformed authentication token"
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			response["message"] = "Invalid token"
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		query := "SELECT username FROM users WHERE id=$1"
		err = database.DBConn.QueryRow(query, tk.UserId).Scan(&user.Username)
		if err == sql.ErrNoRows {
			response["message"] = "User might be deleted from the database"
			w.WriteHeader(http.StatusExpectationFailed)
			json.NewEncoder(w).Encode(response)
			return
		}
		ctx := context.WithValue(r.Context(), UserIdKey{}, user.Username)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	log.Println("------------------------------------------")
	log.Println("App is running on http://127.0.0.1:8000/")
	handleRequests()
}
