package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) (interface{}, *RequestError) {
	return &struct{
		Message string
		Version string
	}{
		Message: "Mahjong Game API",
		Version: "0.1",
	}, nil
}

func (s *Server) handleNew(w http.ResponseWriter, r *http.Request) (interface{}, *RequestError) {
	id := s.Storage.StartNew()
	w.WriteHeader(http.StatusCreated)

	return &struct{
		Message string
		Location string
	}{
		Message: "Game created",
		Location: fmt.Sprintf("%s/game/%d", s.GetDomain(true), id),
	}, nil
}

func (s *Server) handleGame(w http.ResponseWriter, r *http.Request) (interface{}, *RequestError) {
	strId := strings.TrimPrefix(r.URL.Path, "/game/")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return nil, &RequestError{
			Message:    "invalid game id",
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		}
	}

	game, err := s.Storage.Get(uint64(id))
	if err != nil {
		return nil, &RequestError{
			Message: "game does not exist",
			StatusCode: http.StatusNotFound,
			Error: err.Error(),
		}
	}

	return game, nil
}