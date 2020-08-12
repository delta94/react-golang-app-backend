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
	FileAvatar	string		`db:"fileAvatar" json:"fileAvatar"`
	AvatarURL	string		`db:"avatarUrl" json:"avatarUrl`
	CreatedAt	time.Time	`db:"createdAt" json:"createdAt"`
	ModifiedAt	time.Time	`db:"modifiedAt" json:"modifiedAt"`
}