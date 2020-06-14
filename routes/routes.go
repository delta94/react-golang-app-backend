package routes

import (
	"github.com/gorilla/mux"

	a "github.com/marceloOliveira/siteGolang/service/auth"
	p "github.com/marceloOliveira/siteGolang/service/products"
	m "github.com/marceloOliveira/siteGolang/service/user"
)

//Middleware create router and import routes
func Middleware() (r *mux.Router) {
	router := mux.NewRouter()

	router.HandleFunc("/users", m.SelectListOfUser).Methods("GET")
	router.HandleFunc("/users/{id}", m.SelectUser).Methods("GET")
	router.HandleFunc("/users/add", m.InsertUser).Methods("POST")
	router.HandleFunc("/users/update/{id}", m.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/delete/{id}", m.DeleteUser).Methods("DELETE")
	router.HandleFunc("/login", a.AutenticationJWT).Methods("POST")
	router.HandleFunc("/signup", a.SignUp).Methods("POST")
	router.HandleFunc("/signup/usernameList", a.ListUsername).Methods("GET")
	router.HandleFunc("/product", p.SelectProductList).Methods("GET")
	router.HandleFunc("/product/{id}", p.SelectProduct).Methods("GET")
	router.HandleFunc("/product/add", p.InsertProduct).Methods("POST")
	router.HandleFunc("/product/update/{id}", p.UpdateProduct).Methods("PUT")
	router.HandleFunc("/product/delete/{id}", p.DeleteProduct).Methods("DELETE")

	return router
}