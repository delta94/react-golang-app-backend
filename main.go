package main

import (
	"log"
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
	addr := os.Getenv("ADDRESS")
	c.CreateConnection(dbString)
	router := r.Middleware()
	handler := cors.AllowAll().Handler(router)

	log.Fatal(http.ListenAndServe(addr + ":" + port, handler))
}
