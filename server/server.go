package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/roelofruis/mahjong-learn/mahjong"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	Host   string
	Port   string
	Paths  *Paths
	Router *mux.Router

	Games *GameStorage
}

func NewPaths() *Paths {
	return &Paths{
		Index: "/",
		New:   "/new",
		Game:  "/game/",
	}
}

type Paths struct {
	Index string
	New   string
	Game  string
}

type Response struct {
	Data       interface{}
	StatusCode int
	Error      error
}

func (s *Server) GetDomain(includeScheme bool) string {
	if includeScheme {
		return fmt.Sprintf("http://%s:%s", s.Host, s.Port)
	}

	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

type RequestHandler func(*http.Request) *Response

func (s *Server) Routes() {
	s.Router.HandleFunc("/", s.asJsonResponse(s.handleIndex))
	s.Router.HandleFunc("/new", s.asJsonResponse(s.handleNew))
	s.Router.HandleFunc("/game/{id:[0-9]+}", s.asJsonResponse(s.withGame(s.handleDisplayGame))).Methods("GET")
	s.Router.HandleFunc("/game/{id:[0-9]+}", s.asJsonResponse(s.withValidForm(s.withGame(s.handleActions)))).Methods("POST")
	s.Router.HandleFunc("/game/{id:[0-9]+}/player/{seat:[0-9]+}", s.asJsonResponse(s.withGame(s.handleDisplayPlayer))).Methods("GET")
	s.Router.NotFoundHandler = s.asJsonResponse(s.notFoundHandler)
}

func (s *Server) notFoundHandler(_ *http.Request) *Response {
	return &Response{
		StatusCode: http.StatusNotFound,
		Error:      fmt.Errorf("not found"),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *Server) asJsonResponse(f RequestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := f(r)

		if response.StatusCode == 0 {
			http.Error(w, fmt.Sprint("handler returned incomplete Response: status code was not set"), http.StatusInternalServerError)
			return
		}

		var data = response.Data

		if response.Error != nil {
			log.Printf("Handler returned error: %s", response.Error.Error())
			data = struct {
				Error      string `json:"error"`
				StatusCode int    `json:"status_code"`
			}{
				Error:      response.Error.Error(),
				StatusCode: response.StatusCode,
			}
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)

		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to encode data: %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) withGame(f func(r *http.Request, id uint64, game *mahjong.Game) *Response) RequestHandler {
	return func(r *http.Request) *Response {
		vars := mux.Vars(r)
		id, err := intVar(vars, "id")
		if err != nil {
			return &Response{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			}
		}

		g, err := s.Games.Get(uint64(id))
		if err != nil {
			return &Response{
				StatusCode: http.StatusNotFound,
				Error:      fmt.Errorf("no game with id [%d]", id),
			}
		}

		return f(r, uint64(id), g)
	}
}

func (s *Server) withValidForm(f RequestHandler) RequestHandler {
	return func(r *http.Request) *Response {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				return &Response{
					StatusCode: http.StatusBadRequest,
					Error:      fmt.Errorf("unable to parse post data: %s", err.Error()),
				}
			}
		}

		return f(r)
	}
}

func intVar(vars map[string]string, name string) (int, error) {
	strId, has := vars[name]
	if !has {
		return 0, fmt.Errorf("no %s given", name)
	}

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s value [%s]", name, strId)
	}
	return int(id), nil
}
