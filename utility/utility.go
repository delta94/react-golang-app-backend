package utility

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

//JWT key
var jwtKey = []byte("secretKeyApp")

//ResponseStruct struct message with token
type ResponseStruct struct {
	Message		string
	Token		string
	Data		interface{}
}

//Response struct
type Response struct {
	Message		string
}

//Error struct
type Error struct {
	ErrorMsg	string
	Error		error
}

//ResponseJSON struct
type ResponseJSON struct {
	Message		string
	Data		interface{}
}

//ResponseWithJSON to the http body
func ResponseWithJSON(msg string, data interface{}) (resp []byte) {
	msgResp := ResponseJSON{msg, data}
	response, err := json.Marshal(msgResp)
	if err != nil {
		log.Print("Error in marshall", err)
		return
	}
	return response
}

//SucessResponse to http body
func SucessResponse(msg string) (resp []byte) {
	msgResp := Response{msg}
	response, err := json.Marshal(msgResp)
	if err != nil {
		log.Print("Error in marshall", err)
		return
	}
	return response
}

//ErrorResponse to http body
func ErrorResponse(msg string, err error) (resp []byte) {
	errorMsg := Error{msg, err}
	response, err := json.Marshal(errorMsg)
	if err != nil {
		log.Print("Error in marshall", err)
		return
	}
	return response
}

//ResponseWithToken response body with token
func ResponseWithToken(msg string, token string, user interface{}) (resp []byte) {
	msgToken := ResponseStruct{msg, token, user}
	response, err := json.Marshal(msgToken)
	if err != nil {
		log.Print("Error in marshall", err)
		return
	}
	return response
}

//VerifyToken verification of JWT token
func VerifyToken(r *http.Request) (status bool) {
	auth := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	if auth == "" {
		return false
	}
	token, _ := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("There was an error")
        }
        return jwtKey, nil
	})

	return token.Valid
}