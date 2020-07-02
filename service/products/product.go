package product

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"

	u "github.com/marceloOliveira/siteGolang/models"
	d "github.com/marceloOliveira/siteGolang/server"
	a "github.com/marceloOliveira/siteGolang/utility"
)

//SelectProductList Fetch product from database
func SelectProductList(w http.ResponseWriter, r *http.Request)  {
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

	selectDB, err := db.Query("SELECT * FROM products")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	defer selectDB.Close()
	var res = []u.Product{}

	for selectDB.Next() {
		var product u.Product

		err = selectDB.Scan(&product.ProductID, &product.Name, &product.Value, &product.Info, &product.CategoryID)
		if err != nil {
			response := a.ErrorResponse("Error in select", err)
			w.WriteHeader(500)
			w.Write(response)
			return
		}

		res = append(res, product)
	}

	response := a.ResponseWithJSON("Sucess in select", res)
	w.WriteHeader(200)
	w.Write(response)
	db.Close()
}

//SelectProduct specific product
func SelectProduct(w http.ResponseWriter, r *http.Request) {
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

	var product u.Product
	error := json.NewDecoder(r.Body).Decode(&product)
	if error != nil {
		response := a.ErrorResponse("Error in body fields", error)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	selectDB, err := db.Query("SELECT * FROM products WHERE productID = ?", product.ProductID)
	if err != nil {
		response := a.ErrorResponse("Error in query", error)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	defer db.Close()

	var res = []u.Product{}
	for selectDB.Next() {
		var product u.Product
		err = selectDB.Scan(&product.ProductID, &product.Name, &product.Value, &product.Info, &product.CategoryID)
		if err != nil {
			response := a.ErrorResponse("Error in select", error)
			w.WriteHeader(500)
			w.Write(response)
			return
		}

		res = append(res, product)
	}

	response := a.ResponseWithJSON("Sucess in select", res)
	w.WriteHeader(200)
	w.Write(response)
}

//InsertProduct in database
func InsertProduct(w http.ResponseWriter, r *http.Request)  {
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

	stmt, err := db.Prepare("INSERT INTO products(productID, productName, productValue, productInfo, categoryID) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	var product u.Product
	error := json.NewDecoder(r.Body).Decode(&product)
	if error != nil {
		response := a.ErrorResponse("Error in body fields", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	genIDProduct := ksuid.New()
	product.ProductID = genIDProduct
	_, err = stmt.Exec(product.ProductID, product.Name, product.Value, product.Info, product.CategoryID)
	if err != nil {
		response := a.ErrorResponse("Error in insert", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	defer db.Close()

	response := a.SucessResponse("New product add in database")
	w.WriteHeader(200)
	w.Write(response)
}

//DeleteProduct from database
func DeleteProduct(w http.ResponseWriter, r *http.Request)  {
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

	stmt, err := db.Prepare("DELETE FROM products WHERE productID = ?")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	var product u.Product
	error := json.NewDecoder(r.Body).Decode(&product)
	if error != nil {
		response := a.ErrorResponse("Error in fields", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	_, err = stmt.Exec(product.ProductID)
	if err != nil {
		response := a.ErrorResponse("Error in delete", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	response := a.SucessResponse("Product deleted from database")
	w.WriteHeader(200)
	w.Write(response)

	defer db.Close()
}

//UpdateProduct in database
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	stmt, err := db.Prepare("UPDATE products SET productID = ?, productName = ?, productValue = ?, productInfo = ?, categoryID = ? WHERE productID = ?")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	var product u.Product
	error := json.NewDecoder(r.Body).Decode(&product)
	if error != nil {
		response := a.ErrorResponse("Error in fields", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	_, err = stmt.Exec(product.ProductID, product.Name, product.Value, product.Info, product.CategoryID, product.ProductID)
	if err != nil {
		response := a.ErrorResponse("Error in update", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	defer db.Close()

	response := a.SucessResponse("Product updated in the database")
	w.WriteHeader(200)
	w.Write(response)
}