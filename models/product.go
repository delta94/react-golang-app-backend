package models

import (
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/segmentio/ksuid"
)

//Category struct
type Category struct {
	CategoryID		int		`db:"categoryID"`
	Name			string	`db:"categoryName"`
}

//Product struct
type Product struct {
	ProductID		ksuid.KSUID			`db:"productID" json:"productID"`
	Name			string				`db:"productName" json:"Name"`
	Value			float64				`db:"productValue" json:"Value"`
	Info			types.JSONText		`db:"productInfo" json:"Info"`
	FileAvatar		string				`db:"fileAvatar" json:"fileAvatar"`
	AvatarURL		string				`db:"avatarUrl" json:"avatarUrl`
	CategoryID		int					`db:"categoryID" json:"CategoryID"`
	CreatedAt		time.Time			`db:"createdAt" json:"createdAt"`
	ModifiedAt		time.Time			`db:"modifiedAt" json:"modifiedAt"`
}