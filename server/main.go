package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := &Server{
		Host: "localhost",
		Port: port,
		Router: http.ServeMux{},

		Games: NewGameStorage(),
	}

	server.Routes()

	log.Printf("mahjong API")
	log.Printf("server starting on %s", server.GetDomain(true))
	err := http.ListenAndServe(server.GetDomain(false), server)
	if err != nil {
		log.Fatal(err)
	}
}
