package models

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
