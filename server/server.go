package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	Host string
	Port string
	Router http.ServeMux

	Games *GameStorage
}

type Response struct {
	Data interface{}
	StatusCode int
	Error error
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
	s.Router.HandleFunc("/show/", s.asJsonResponse(s.withGame(s.handleShow)))
	s.Router.HandleFunc("/advance/", s.asJsonResponse(s.withGame(s.handleAdvance)))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) asJsonResponse(f func(r *http.Request) *Response) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		response := f(r)

		var data = response.Data

		if response.Error != nil {
			log.Printf("Handler returned error: %s", response.Error.Error())
			data = response.Error
		}

		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to encode data: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(response.StatusCode)
	}
}

func (s *Server) withGame(f func (r *http.Request, game Game) *Response) func(r *http.Request) *Response {
	return func (r *http.Request) *Response {
		strId := strings.TrimPrefix(r.URL.Path, "/game/")
		id, err := strconv.ParseInt(strId, 10, 64)
		if err != nil {
			return &Response{
				StatusCode: http.StatusBadRequest,
				Error: fmt.Errorf("invalid game id [%s]", strId),
			}
		}

		game, err := s.Games.Get(uint64(id))
		if err != nil {
			return &Response{
				StatusCode: http.StatusNotFound,
				Error: fmt.Errorf("no game with id [%s]", strId),
			}
		}

		return f(r, game)
	}
}