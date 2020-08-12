package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	u "github.com/marceloOliveira/siteGolang/models"
	d "github.com/marceloOliveira/siteGolang/server"
	a "github.com/marceloOliveira/siteGolang/utility"
)

//JWT key
func getJWTkey() ([]byte) {
	godotenv.Load(".env")
	jwtKey := []byte(os.Getenv("JWT_KEY"))
	return jwtKey
}

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
		var createdAt []uint8
		var modifiedAt []uint8
		err = query.Scan(&user.UserID, &user.Username, &user.Password, &user.Fullname, &user.FileAvatar, &user.AvatarURL, &createdAt, &modifiedAt)
		if err != nil {
			response := a.ErrorResponse("Error in select", err)
			fmt.Println(err)
			w.WriteHeader(500)
			w.Write(response)
			return
		}
		user.CreatedAt, _ = a.ConvertToTime(createdAt)
		user.ModifiedAt, _ = a.ConvertToTime(modifiedAt)
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
	
	expirationTime := time.Now().Add(5000 * time.Minute)
	claims := &u.Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := getJWTkey()
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
	db.Close()
}

//SignUp function to add user to database withou jwt
func SignUp(w http.ResponseWriter, r *http.Request)  {
	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	db := d.CreateConnection(dbString)
	w.Header().Set("Content-Type", "multipart/form-data")

	maxSize := int64(3 * 1024 * 1024)
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		response := a.ErrorResponse("Failed to parse multipart form", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	useridGenerate := uuid.Must(uuid.NewRandom())
	var user u.User
	user.UserID = useridGenerate
	user.Username = r.FormValue("username")
	user.Fullname = r.FormValue("fullname")

	stmt, err := db.Prepare("INSERT INTO users(userID, username, password, fullName, fileAvatar, avatarUrl, createdAt, modifiedAt) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	
	bytePassword := []byte(r.FormValue("password"))
	hashPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
    if err != nil {
		response := a.ErrorResponse("Error generating hash password", err)
		w.WriteHeader(500)
		w.Write(response)
        return
	}
	user.Password = string(hashPassword)
	user.CreatedAt = time.Now()

	avatarFile, fileHeader, err := r.FormFile("avatar")
	if err != nil {
		response := a.ErrorResponse("Failed to get upload file", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	defer avatarFile.Close()

	secretKey := os.Getenv("AWS_SECRET_KEY")
	secretID := os.Getenv("AWS_SECRET_ID")
	region := os.Getenv("AWS_REGION")
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(secretID, secretKey, ""),
	})
	if err != nil {
		response := a.ErrorResponse("Failed to set session on AWS S3", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	fileName, imageURL, err := a.UploadImageToS3(session, avatarFile, fileHeader)
	if err != nil {
		response := a.ErrorResponse("Failed to upload Image to S3", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	user.FileAvatar = fileName
	user.AvatarURL = imageURL

	_, err = stmt.Exec(user.UserID, user.Username, user.Password, user.Fullname, user.FileAvatar, user.AvatarURL, user.CreatedAt, user.ModifiedAt)
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
	db.Close()
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
	var data []string

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

	for i := 0; i < len(res); i++ {
		data = append(data, res[i].Username)
	}

	response := a.ResponseWithJSON("Success in select from database", data)
	w.WriteHeader(200)
	w.Write(response)
	db.Close()
}