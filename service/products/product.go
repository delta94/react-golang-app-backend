package product

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

		err = selectDB.Scan(&product.ProductID, &product.Name, &product.Value, &product.Info, &product.CategoryID, &product.FileAvatar, &product.AvatarURL)
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
	prodID := mux.Vars(r)
	if(prodID["id"] == "") {
		var err error
		response := a.ErrorResponse("Missing field product id", err)
		w.WriteHeader(400)
		w.Write(response)
		return
	}
	
	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	db := d.CreateConnection(dbString)
	w.Header().Set("Content-Type", "application/json")

	selectDB, err := db.Query("SELECT * FROM products WHERE productID = ?", prodID["id"])
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}
	defer db.Close()

	var res = []u.Product{}
	for selectDB.Next() {
		var product u.Product
		err = selectDB.Scan(&product.ProductID, &product.Name, &product.Value, &product.Info, &product.CategoryID, &product.FileAvatar, &product.AvatarURL)
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

	maxSize := int64(3 * 1024 * 1024)
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		response := a.ErrorResponse("Failed to parse multipart form", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	stmt, err := db.Prepare("INSERT INTO products(productID, productName, productValue, productInfo, categoryID) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	var product u.Product
	genIDProduct := ksuid.New()
	product.ProductID = genIDProduct
	product.CreatedAt = time.Now()
	product.Name = r.FormValue("name")
	product.Value, _ = strconv.ParseFloat(r.FormValue("value"), 64)
	product.CategoryID, _ = strconv.Atoi(r.FormValue("category"))
	product.Info = []byte(r.FormValue("info"))


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
	prodID := mux.Vars(r)

	stmt, err := db.Prepare("DELETE FROM products WHERE productID = ?")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	_, err = stmt.Exec(prodID["id"])
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