package models

import "github.com/golang-jwt/jwt"

type Token struct {
	UserId uint
	jwt.StandardClaims
}
