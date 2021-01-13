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
	Paths *Paths
	Router http.ServeMux

	Games *GameStorage
}

func NewPaths() *Paths {
	return &Paths{
		Index: "/",
		New: "/new",
		Show: "/show/",
		Advance: "/advance/",
	}
}

type Paths struct {
	Index string
	New string
	Show string
	Advance string
}

type Response struct {
	Data interface{}
	StatusCode int
	Error error
}

type ErrorMessage struct {
	Error string
	StatusCode int
}

func (s *Server) GetDomain(includeScheme bool) string {
	if includeScheme {
		return fmt.Sprintf("http://%s:%s", s.Host, s.Port)
	}

	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

func (s *Server) Routes() {
	s.Router.HandleFunc(s.Paths.Index, s.asJsonResponse(s.handleIndex))
	s.Router.HandleFunc(s.Paths.New, s.asJsonResponse(s.handleNew))
	s.Router.HandleFunc(s.Paths.Show, s.asJsonResponse(s.withGame(s.handleShow)))
	s.Router.HandleFunc(s.Paths.Advance, s.asJsonResponse(s.withGame(s.handleAdvance)))
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
			data = &ErrorMessage{StatusCode: response.StatusCode, Error: response.Error.Error()}
		}

		w.WriteHeader(response.StatusCode)

		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to encode data: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) withGame(f func (r *http.Request, game Game) *Response) func(r *http.Request) *Response {
	return func (r *http.Request) *Response {
		parts := strings.Split(r.URL.Path, "/")
		var strId string
		if len(parts) == 2 {
			strId = parts[1]
		} else if len(parts) == 3 {
			strId = parts[2]
		} else {
			return &Response{
				StatusCode: http.StatusBadRequest,
				Error: fmt.Errorf("unable to determine id"),
			}
		}

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