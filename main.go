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
	portEnv := os.Getenv("PORT")
	var port string
	if portEnv != "" {
		port = portEnv
	} else {
		port = ":4000"
	}
	c.CreateConnection(dbString)
	router := r.Middleware()
	handler := cors.AllowAll().Handler(router)

	fmt.Println("Server start on localhost" + port)
	http.ListenAndServe(port, handler)
}
