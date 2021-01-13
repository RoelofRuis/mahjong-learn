package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) handleIndex(_ *http.Request) *Response {
	return &Response{
		StatusCode: http.StatusFound,
		Data: &struct {
			Message string
			Version string
		}{
			Message: "Mahjong Game API",
			Version: "0.1",
		},
	}
}

func (s *Server) handleNew(r *http.Request) *Response {
	id := s.Games.StartNew()

	return &Response{
		StatusCode: http.StatusCreated,
		Data: &struct{
			Message string
			Location string
		}{
			Message: "Game created",
			Location: fmt.Sprintf("%s/show/%d", s.GetDomain(true), id),
		},
	}
}

func (s *Server) handleShow(r *http.Request) *Response {
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

	return &Response{
		StatusCode: http.StatusFound,
		Data: game,
	}
}

func (s *Server) handleAdvance(r *http.Request) *Response {
	// TODO: implement
	return &Response{}
}
