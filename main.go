package main

import (
	"exp/routes"
	_"exp/docs"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	e := godotenv.Load()
	if e != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(e)
	port := os.Getenv("PORT")
	http.Handle("/", routes.Routers())

	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}