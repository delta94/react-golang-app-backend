package models

import (
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
	Info			*types.JSONText		`db:"productInfo" json:"Info"`
	CategoryID		int					`db:"categoryID" json:"CategoryID"`
}