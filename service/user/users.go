package user

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/uuid"
	u "github.com/marceloOliveira/siteGolang/models"
	d "github.com/marceloOliveira/siteGolang/server"
	a "github.com/marceloOliveira/siteGolang/utility"
	"golang.org/x/crypto/bcrypt"
)

//SelectListOfUser in database
func SelectListOfUser(w http.ResponseWriter, r *http.Request) {
	auth := a.VerifyToken(r)
	if !auth {
		var err error
		response := a.ErrorResponse("Unauthorized", err)
		w.WriteHeader(401)
		w.Write(response)
		return
	}
	
	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	db := d.CreateConnection(dbString)

	selectDB, err := db.Query("SELECT * FROM users")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	defer selectDB.Close()
	var res = []u.User{}

	for selectDB.Next() {
		var user u.User
		var createdAt []uint8
		var modifiedAt []uint8
		err = selectDB.Scan(&user.UserID, &user.Username, &user.Password, &user.Fullname, &user.FileAvatar, &user.AvatarURL, &createdAt, &modifiedAt)
		if err != nil {
			response := a.ErrorResponse("Error in select", err)
			fmt.Println(err)
			w.WriteHeader(500)
			w.Write(response)
			return
		}
		user.CreatedAt, _ = a.ConvertToTime(createdAt)
		user.ModifiedAt, _ = a.ConvertToTime(modifiedAt)
		res = append(res, user)
	}
	defer db.Close()

	response := a.ResponseWithJSON("Success in select from database", res)
	w.WriteHeader(200)
	w.Write(response)
}

//SelectUser specific in database
func SelectUser(w http.ResponseWriter, r *http.Request)  {
	auth := a.VerifyToken(r)
	if !auth {
		var err error
		response := a.ErrorResponse("Unauthorized", err)
		w.WriteHeader(401)
		w.Write(response)
		return
	}
	vars := mux.Vars(r)
	if(vars["id"] == "") {
		var err error
		response := a.ErrorResponse("Missing field user id", err)
		w.WriteHeader(400)
		w.Write(response)
		return
	}

	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	db := d.CreateConnection(dbString)

	selectDB, err := db.Query("SELECT * FROM users WHERE userID = ?", vars["id"])
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	defer db.Close()

	var res []u.User
	for selectDB.Next() {
		var user u.User
		var createdAt []uint8
		var modifiedAt []uint8
		err = selectDB.Scan(&user.UserID, &user.Username, &user.Password, &user.Fullname, &user.FileAvatar, &user.AvatarURL, &createdAt, &modifiedAt)
		if err != nil {
			response := a.ErrorResponse("Error in select", err)
			w.WriteHeader(500)
			w.Write(response)
			return
		}
		user.CreatedAt, _ = a.ConvertToTime(createdAt)
		user.ModifiedAt, _ = a.ConvertToTime(modifiedAt)
		res = append(res, user)
	}

	response := a.ResponseWithJSON("Success in select from database", res)
	w.WriteHeader(200)
	w.Write(response)
}

//InsertUser in database
func InsertUser(w http.ResponseWriter, r *http.Request) {
	auth := a.VerifyToken(r)
	if !auth {
		var err error
		response := a.ErrorResponse("Unauthorized", err)
		w.WriteHeader(401)
		w.Write(response)
		return
	}

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
	var user u.User
	useridGenerate := uuid.Must(uuid.NewRandom())
	user.UserID = useridGenerate
	user.Username = r.FormValue("username")
	user.Fullname = r.FormValue("fullname")
	user.CreatedAt = time.Now()
	bytePassword := []byte(r.FormValue("password"))
	hashPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
    if err != nil {
		response := a.ErrorResponse("Error generating hash password", err)
		w.WriteHeader(500)
		w.Write(response)
        return
	}
	user.Password = string(hashPassword)

	if r.FormValue("hasAvatar") == "true" {
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
	} else if r.FormValue("hasAvatar") == "false" {
		user.FileAvatar = ""
		user.AvatarURL = ""
	}

	stmt, err := db.Prepare("INSERT INTO users(userID, username, password, fullName, fileAvatar, avatarUrl, createdAt, modifiedAt) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

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
}

//UpdateUser in database
func UpdateUser(w http.ResponseWriter, r *http.Request)  {
	auth := a.VerifyToken(r)
	if !auth {
		var err error
		response := a.ErrorResponse("Unauthorized", err)
		w.WriteHeader(401)
		w.Write(response)
		return
	}

	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	db := d.CreateConnection(dbString)
	w.Header().Set("Content-Type", "multipart/form-data")

	stmt, err := db.Prepare("UPDATE users SET username = ?, password = ?, fullName = ?, fileAvatar = ?, avatarUrl = ?, modifiedAt = ? WHERE userID = ?")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	vars := mux.Vars(r)
	var user u.User
	maxSize := int64(3 * 1024 * 1024)
	errorParse := r.ParseMultipartForm(maxSize)
	if errorParse != nil {
		response := a.ErrorResponse("Failed to parse multipart form", errorParse)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	user.Username = r.FormValue("username")
	user.Fullname = r.FormValue("fullname")
	user.ModifiedAt = time.Now()
	if r.FormValue("hasAvatar") == "true" {
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
	} else if r.FormValue("hasAvatar") == "false" {
		user.FileAvatar = ""
		user.AvatarURL = ""
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

	_, err = stmt.Exec(user.Username, user.Password, user.Fullname, user.FileAvatar, user.AvatarURL, user.ModifiedAt, vars["id"])
	if err != nil {
		response := a.ErrorResponse("Error in update", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	defer db.Close()

	response := a.SucessResponse("User updated in the database")
	w.WriteHeader(200)
	w.Write(response)
}

//DeleteUser in database
func DeleteUser(w http.ResponseWriter, r *http.Request)  {
	auth := a.VerifyToken(r)
	if !auth {
		var err error
		response := a.ErrorResponse("Unauthorized", err)
		w.WriteHeader(401)
		w.Write(response)
		return
	}
	userID := mux.Vars(r)
	if(userID["id"] == "") {
		var err error
		response := a.ErrorResponse("Missing field user id", err)
		w.WriteHeader(400)
		w.Write(response)
		return
	}
	
	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	db := d.CreateConnection(dbString)

	stmt, err := db.Prepare("DELETE FROM users WHERE userID = ?")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	_, err = stmt.Exec(userID["id"])
	if err != nil {
		response := a.ErrorResponse("Error in delete", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	response := a.SucessResponse("User deleted from database")
	w.WriteHeader(200)
	w.Write(response)

	defer db.Close()
}