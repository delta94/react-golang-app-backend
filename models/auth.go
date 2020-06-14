package models

import (
	"github.com/dgrijalva/jwt-go"
)

//Claims struct of encoded JWT
type Claims struct {
	Username	string	`json:"username"`
	jwt.StandardClaims
}