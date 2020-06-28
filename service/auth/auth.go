package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	u "github.com/marceloOliveira/siteGolang/models"
	d "github.com/marceloOliveira/siteGolang/server"
	a "github.com/marceloOliveira/siteGolang/utility"
)

//JWT key
var jwtKey = []byte("secretKeyApp")

//AutenticationJWT generate JWT token
func AutenticationJWT(w http.ResponseWriter, r *http.Request)  {
	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	db := d.CreateConnection(dbString)
	w.Header().Set("Content-Type", "application/json")

	var userBody u.User
	error := json.NewDecoder(r.Body).Decode(&userBody)
	if error != nil {
		response := a.ErrorResponse("Error in body fields", error)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	query, err := db.Query("SELECT * FROM users WHERE username = ?", userBody.Username)
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	defer query.Close()

	var user u.User
	for query.Next() {
		err = query.Scan(&user.UserID, &user.Username, &user.Password, &user.Fullname, &user.Avatar)
		if err != nil {
			response := a.ErrorResponse("Error in select", err)
			w.WriteHeader(500)
			w.Write(response)
			return
		}
	}

	bytePassword := []byte(user.Password)
	bodyPassword := []byte(userBody.Password)
	errorHash := bcrypt.CompareHashAndPassword(bytePassword, bodyPassword)
	if errorHash != nil {
		response := a.ErrorResponse("Invalid Password!", err)
		w.WriteHeader(400)
		w.Write(response)
		return
	}
	
	expirationTime := time.Now().Add(1440 * time.Minute)
	claims := &u.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		response := a.ErrorResponse("Error generate jwt token", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	response := a.ResponseWithToken("Sucess token generated", tokenString, user)
	w.WriteHeader(200)
	w.Write(response)
}

//SignUp function to add user to database withou jwt
func SignUp(w http.ResponseWriter, r *http.Request)  {
	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	log.Print("string ", dbString)
	db := d.CreateConnection(dbString)
	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("INSERT INTO users(userID, username, password, fullName, avatar) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	
	useridGenerate := uuid.Must(uuid.NewRandom())
	var user u.User
	error := json.NewDecoder(r.Body).Decode(&user)
	if error != nil {
		response := a.ErrorResponse("Error in body fields", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	user.UserID = useridGenerate
	bytePassword := []byte(user.Password)
	hashPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
    if err != nil {
		response := a.ErrorResponse("Error generating hash password", err)
		w.WriteHeader(500)
		w.Write(response)
        return
	}
	user.Password = string(hashPassword)

	_, err = stmt.Exec(user.UserID, user.Username, user.Password, user.Fullname, user.Avatar)
	if err != nil {
		response := a.ErrorResponse("Error in insert", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	defer db.Close()

	response := a.SucessResponse("New user add in database")
	w.WriteHeader(200)
	w.Write(response)
}

//ListUsername get username from the database
func ListUsername(w http.ResponseWriter, r *http.Request)  {
	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	db := d.CreateConnection(dbString)

	query, err := db.Query("SELECT username FROM users")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	defer query.Close()
	var res = []u.User{}

	for query.Next() {
		var user u.User

		err = query.Scan(&user.Username)
		if err != nil {
			response := a.ErrorResponse("Error in select", err)
			w.WriteHeader(500)
			w.Write(response)
			return
		}

		res = append(res, user)
	}

	response := a.ResponseWithJSON("Success in select from database", res)
	w.WriteHeader(200)
	w.Write(response)
}