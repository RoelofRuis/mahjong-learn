package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/roelofruis/mahjong-learn/mahjong"
	"github.com/roelofruis/mahjong-learn/state"
	"github.com/roelofruis/mahjong-learn/view"
	"net/http"
	"strconv"
)

func (s *Server) handleIndex(_ *http.Request) *Response {
	return &Response{
		StatusCode: http.StatusOK,
		Data: &struct {
			Message      string `json:"message"`
			Version      string `json:"version"`
			GamesStarted int    `json:"games_started"`
			NewGame      string `json:"new_game"`
		}{
			Message:      "Mahjong Game API",
			Version:      "0.1",
			GamesStarted: int(*s.Games.lastIndex),
			NewGame:      fmt.Sprintf("%s/new", s.GetDomain(true)),
		},
	}
}

func (s *Server) handleNew(r *http.Request) *Response {
	id, err := s.Games.StartNew()
	if err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	return &Response{
		StatusCode: http.StatusCreated,
		Data: &struct {
			Message  string `json:"message"`
			Id       uint64 `json:"id"`
			Location string `json:"location"`
		}{
			Message:  "Game created",
			Id:       id,
			Location: fmt.Sprintf("%s/game/%d", s.GetDomain(true), id),
		},
	}
}

func (s *Server) handleDisplayGame(r *http.Request, game *mahjong.Game, _ uint64) *Response {
	return &Response{
		StatusCode: http.StatusOK,
		Data:       view.ViewGame(game),
	}
}

func (s *Server) handleDisplayPlayer(r *http.Request, game *mahjong.Game, _ uint64) *Response {
	seat, err := intVar(mux.Vars(r), "seat")
	if err != nil {
		return &Response{
			StatusCode: http.StatusBadRequest,
			Error:      err,
		}
	}
	fmt.Printf("%+v", seat)
	if seat < 0 || seat > 4 {
		return &Response{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("player should be between 0 and 3 inclusive"),
		}
	}
	return &Response{
		StatusCode: http.StatusOK,
		Data:       view.ViewPlayer(game, state.Seat(seat)),
	}
}

func (s *Server) handleActions(r *http.Request, game *mahjong.Game, id uint64) *Response {
	actionMap := make(map[state.Seat]int)
	for i, playerKey := range []string{"1", "2", "3", "4"} {
		playerAction, err := strconv.ParseInt(r.PostForm.Get(playerKey), 10, 64)
		if err == nil {
			actionMap[state.Seat(i)] = int(playerAction)
		}
	}

	err := game.StateMachine.Transition(actionMap)
	if err != nil {
		if _, ok := err.(*state.IncorrectActionError); ok {
			return &Response{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			}
		}

		return &Response{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	return &Response{
		StatusCode: http.StatusAccepted,
		Data: &struct {
			Message  string `json:"message"`
			Id       uint64 `json:"id"`
			Location string `json:"location"`
		}{
			Message:  "actions executed",
			Location: fmt.Sprintf("%s/game/%d", s.GetDomain(true), id),
		},
	}
}
