package product

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/segmentio/ksuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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

		err = selectDB.Scan(&product.ProductID, &product.Name, &product.Value, &product.Info, &product.CategoryID, &product.FileAvatar, &product.AvatarURL, &product.CreatedAt, &product.ModifiedAt)
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
		err = selectDB.Scan(&product.ProductID, &product.Name, &product.Value, &product.Info, &product.CategoryID, &product.FileAvatar, &product.AvatarURL, &product.CreatedAt, &product.ModifiedAt)
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
	w.Header().Set("Content-Type", "multipart/form-data")

	maxSize := int64(3 * 1024 * 1024)
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		response := a.ErrorResponse("Failed to parse multipart form", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	stmt, err := db.Prepare("INSERT INTO products(productID, productName, productValue, fileAvatar, avatarUrl, productInfo, categoryID, createdAt, modifiedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	if(r.FormValue("name") == "" || r.FormValue("value") == "" || r.FormValue("category") == "" || r.FormValue("info") == "") {
		var err error
		response := a.ErrorResponse("Missing fields", err)
		w.WriteHeader(400)
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
		product.FileAvatar = fileName
		product.AvatarURL = imageURL


	_, err = stmt.Exec(product.ProductID, product.Name, product.Value, product.FileAvatar, product.AvatarURL, product.Info, product.CategoryID, product.CreatedAt, product.ModifiedAt)
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
	w.Header().Set("Content-Type", "multipart/form-data")
	
	maxSize := int64(3 * 1024 * 1024)
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		response := a.ErrorResponse("Failed to parse multipart form", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	productID := mux.Vars(r)
	if(productID["id"] == "" || r.FormValue("name") == "" || r.FormValue("value") == "" || r.FormValue("category") == "" || r.FormValue("info") == "") {
		var err error
		response := a.ErrorResponse("Missing fields", err)
		w.WriteHeader(400)
		w.Write(response)
		return
	}

	var product u.Product
	product.ModifiedAt = time.Now()
	product.Name = r.FormValue("name")
	product.Value, _ = strconv.ParseFloat(r.FormValue("value"), 64)
	product.CategoryID, _ = strconv.Atoi(r.FormValue("category"))
	product.Info = []byte(r.FormValue("info"))
	product.CreatedAt = time.Now()
	
	if r.FormValue("hasAvatar") == "true"  {
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
		product.FileAvatar = fileName
		product.AvatarURL = imageURL
	} else if r.FormValue("hasAvatar") == "false" {
		product.FileAvatar = ""
		product.AvatarURL = ""
	}
    
	if r.FormValue("hasAvatar") == "true" {
		stmt, err := db.Prepare("UPDATE products SET productName = ?, productValue = ?, productInfo = ?, fileAvatar = ?, avatarUrl = ?, categoryID = ?, modifiedAt = ? WHERE productID = ?")
		if err != nil {
			response := a.ErrorResponse("Error in query", err)
			w.WriteHeader(500)
			w.Write(response)
			return
		}
		_, err = stmt.Exec(product.Name, product.Value, product.Info, product.FileAvatar, product.AvatarURL, product.CategoryID, product.ModifiedAt, productID["id"])
		if err != nil {
			response := a.ErrorResponse("Error in update", err)
			w.WriteHeader(500)
			w.Write(response)
			return
		}
		defer db.Close()
	} else if r.FormValue("hasAvatar") == "false" {
		stmt, err := db.Prepare("UPDATE products SET productName = ?, productValue = ?, productInfo = ?, categoryID = ?, modifiedAt = ? WHERE productID = ?")
		if err != nil {
			response := a.ErrorResponse("Error in query", err)
			w.WriteHeader(500)
			w.Write(response)
			return
		}
		_, err = stmt.Exec(product.Name, product.Value, product.Info, product.CategoryID, product.ModifiedAt, productID["id"])
		if err != nil {
			response := a.ErrorResponse("Error in update", err)
			w.WriteHeader(500)
			w.Write(response)
			return
		}
		defer db.Close()
	}

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
	
    productID := mux.Vars(r)
	if(productID["id"] == "") {
		var err error
		response := a.ErrorResponse("Missing field product id", err)
		w.WriteHeader(400)
		w.Write(response)
		return
	}
	
	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	db := d.CreateConnection(dbString)

	stmt, err := db.Prepare("DELETE FROM products WHERE productID = ?")
	if err != nil {
		response := a.ErrorResponse("Error in query", err)
		w.WriteHeader(500)
		w.Write(response)
		return
	}

	_, err = stmt.Exec(productID["id"])
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