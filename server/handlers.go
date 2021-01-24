package main

import (
	"fmt"
	"github.com/roelofruis/mahjong-learn/game"
	"net/http"
	"strconv"
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

func (s *Server) handleDisplay(r *http.Request, stateMachine *game.StateMachine) *Response {
	return &Response{
		StatusCode: http.StatusFound,
		Data:       View(stateMachine),
	}
}

func (s *Server) handleTransition(r *http.Request, stateMachine *game.StateMachine) *Response {
	actionMap := make(map[game.Seat]int)
	for i, playerKey := range []string{"1", "2", "3", "4"} {
		playerAction, err := strconv.ParseInt(r.PostForm.Get(playerKey), 10, 64)
		if err == nil {
			actionMap[game.Seat(i)] = int(playerAction)
		}
	}

	err := stateMachine.Transition(actionMap)
	if err != nil {
		return &Response{
			StatusCode: http.StatusBadRequest,
			Error:      err,
		}
	}

	return &Response{
		StatusCode: http.StatusAccepted,
		Data: &struct {
			Message  string
			Id       uint64
			Location string
		}{
			Message:  "StateTransition executed",
			Location: fmt.Sprintf("%s%s%d", s.GetDomain(true), s.Paths.Game, stateMachine.Id()),
		},
	}
}
