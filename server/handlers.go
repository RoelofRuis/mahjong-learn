package main

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/game"
	"net/http"
)

func (s *Server) handleIndex(_ *http.Request) *Response {
	return &Response{
		StatusCode: http.StatusFound,
		Data: &struct {
			Message string
			Version string
			NewGame string
		}{
			Message: "Mahjong Game API",
			Version: "0.1",
			NewGame: fmt.Sprintf("%s%s", s.GetDomain(true), s.Paths.New),
		},
	}
}

func (s *Server) handleNew(r *http.Request) *Response {
	id := s.Games.StartNew()

	return &Response{
		StatusCode: http.StatusCreated,
		Data: &struct {
			Message  string
			Id       uint64
			Location string
		}{
			Message:  "Game created",
			Id:       id,
			Location: fmt.Sprintf("%s%s%d", s.GetDomain(true), s.Paths.Game, id),
		},
	}
}

func (s *Server) handleShow(r *http.Request, stateMachine *game.StateMachine) *Response {
	return &Response{
		StatusCode: http.StatusFound,
		Data:       View(stateMachine),
	}
}

func (s *Server) handleAction(r *http.Request, stateMachine *game.StateMachine) *Response {
	// TODO: implement
	return &Response{}
}
