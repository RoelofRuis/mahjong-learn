package main

import (
	"fmt"
	"net/http"
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

func (s *Server) handleShow(r *http.Request, game Game) *Response {
	return &Response{
		StatusCode: http.StatusFound,
		Data: game,
	}
}

func (s *Server) handleAdvance(r *http.Request, game Game) *Response {
	// TODO: implement
	return &Response{}
}
