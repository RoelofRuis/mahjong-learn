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
			Location string
		}{
			Message: "Mahjong Game API",
			Version: "0.1",
			Location: fmt.Sprintf("%s%s", s.GetDomain(true), s.Paths.New),
		},
	}
}

func (s *Server) handleNew(r *http.Request) *Response {
	id := s.Games.StartNew()

	return &Response{
		StatusCode: http.StatusCreated,
		Data: &struct{
			Message string
			Id uint64
			Location string
		}{
			Message: "Game created",
			Id: id,
			Location: fmt.Sprintf("%s%s%d", s.GetDomain(true), s.Paths.Show, id),
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
