package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := 8000

	server := &Server{
		Router: http.ServeMux{},
	}

	server.Routes()

	log.Printf("mahjong API")
	log.Printf("server started on localhost:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), server)
	if err != nil {
		log.Fatal(err)
	}
}
