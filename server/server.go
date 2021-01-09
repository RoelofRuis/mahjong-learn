package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Router http.ServeMux
}

func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.asJsonResponse(handleIndex))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

type RequestError struct {
	Message string
}

func (s *Server) asJsonResponse(f func(r *http.Request) (interface{}, error)) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		obj, err := f(r)
		if err != nil {
			log.Printf("method returned error: %s", err.Error())
			obj = RequestError{Message: err.Error()}
		}

		err = json.NewEncoder(w).Encode(obj)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to encode data: %s", err.Error()), http.StatusInternalServerError)
		}
	}
}