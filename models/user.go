package models

import (
	"time"

	"github.com/google/uuid"
)

//User struct database
type User struct{
	UserID 		uuid.UUID 	`db:"userID" json:"userID"`
	Username 	string 		`db:"username" json:"username"`
	Password 	string 		`db:"password" json:"password"`
	Fullname	string		`db:"fullName" json:"fullname"`
	Avatar		string		`db:"avatar" json:"avatar"`
	CreatedAt	time.Time	`db:"createdAt" json:"createdAt"`
	ModifiedAt	time.Time	`db:"modifiedAt" json:"modifiedAt"`
}