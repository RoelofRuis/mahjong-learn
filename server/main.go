package main

import (
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := &Server{
		Host:   "localhost",
		Port:   port,
		Paths:  NewPaths(),
		Router: mux.NewRouter(),

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
