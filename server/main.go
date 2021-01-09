package main

import (
	"fmt"
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
		Router: http.ServeMux{},
	}

	server.Routes()

	log.Printf("mahjong API")
	log.Printf("server started on localhost:%s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), server)
	if err != nil {
		log.Fatal(err)
	}
}
