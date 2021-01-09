package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Host string
	Port string
	Router http.ServeMux

	Storage *GameStorage
}

type RequestError struct {
	Message string
	StatusCode int

	Error string
}

func (s *Server) GetDomain(includeScheme bool) string {
	if includeScheme {
		return fmt.Sprintf("http://%s:%s", s.Host, s.Port)
	}

	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.asJsonResponse(s.handleIndex))
	s.Router.HandleFunc("/new", s.asJsonResponse(s.handleNew))
	s.Router.HandleFunc("/game/", s.asJsonResponse(s.handleGame))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) asJsonResponse(f func(w http.ResponseWriter, r *http.Request) (interface{}, *RequestError)) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		obj, requestError := f(w, r)
		if requestError != nil {
			log.Printf("Handler returned request error: %s", requestError.Error)
			w.WriteHeader(requestError.StatusCode)
			obj = requestError
		}

		err := json.NewEncoder(w).Encode(obj)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to encode data: %s", err.Error()), http.StatusInternalServerError)
		}
	}
}