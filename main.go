package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/rs/cors"

	r "github.com/marceloOliveira/siteGolang/routes"
	c "github.com/marceloOliveira/siteGolang/server"
)

func main() {
	godotenv.Load(".env")
	dbString := os.Getenv("DBSTRING")
	c.CreateConnection(dbString)
	router := r.Middleware()
	handler := cors.AllowAll().Handler(router)

	fmt.Println("Server start on localhost:4000")
	http.ListenAndServe(":4000", handler)
}
