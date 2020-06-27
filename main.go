package main

import (
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
	port := os.Getenv("PORT")
	c.CreateConnection(dbString)
	router := r.Middleware()
	handler := cors.AllowAll().Handler(router)

	http.ListenAndServe(port, handler)
}
