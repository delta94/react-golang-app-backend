package user

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"

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
	w.Header().Set("Content-Type", "application/json")

	selectDB, err := db.Query("SELECT * FROM mp14jxypt0bem0vd.users")
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

		err = selectDB.Scan(&user.UserID, &user.Username, &user.Password, &user.Fullname, &user.Avatar)
		if err != nil {
			response := a.ErrorResponse("Error in select", err)
			w.WriteHeader(500)
			w.Write(response)
			return
		}

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

	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	db := d.CreateConnection(dbString)
	w.Header().Set("Content-Type", "application/json")

	var userID u.User
	error := json.NewDecoder(r.Body).Decode(&userID)
	if error != nil {
		response := a.ErrorResponse("Error in body fields", error)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	selectDB, err := db.Query("SELECT * FROM users WHERE userID = ?", userID.UserID)
	if err != nil {
		response := a.ErrorResponse("Error in query", error)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	defer db.Close()

	var res []u.User
	for selectDB.Next() {
		var user u.User
		err = selectDB.Scan(&user.UserID, &user.Username, &user.Password, &user.Fullname, &user.Avatar)
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
	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("UPDATE users SET userID = ?, username = ?, password = ?, fullName = ?, avatar = ? WHERE userID = ?")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	var user u.User
	error := json.NewDecoder(r.Body).Decode(&user)
	if error != nil {
		response := a.ErrorResponse("Error in body fields", error)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	bytePassword := []byte(user.Password)
	hashPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		response := a.ErrorResponse("Error generating hash password", err)
		w.WriteHeader(500)
		w.Write(response)
        return
	}
	user.Password = string(hashPassword)

	_, err = stmt.Exec(user.UserID, user.Username, user.Password, user.Fullname, user.Avatar, user.UserID)
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
	
	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	db := d.CreateConnection(dbString)
	w.Header().Set("Content-Type", "application/json")

	stmt, err := db.Prepare("DELETE FROM users WHERE userID = ?")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	var userid u.User
	error := json.NewDecoder(r.Body).Decode(&userid)
	if error != nil {
		response := a.ErrorResponse("Error in fields", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	_, err = stmt.Exec(userid.UserID)
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