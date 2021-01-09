package main

import (
	"errors"
	"sync"
	"sync/atomic"
)

func NewGameStorage() *GameStorage {
	return &GameStorage{
		gamesLock: sync.RWMutex{},
		games:     make(map[uint64]Game),
		lastIndex: new(uint64),
	}
}

type GameStorage struct {
	gamesLock sync.RWMutex
	games map[uint64]Game

	lastIndex *uint64
}

func (s *GameStorage) Get(id uint64) (Game, error) {
	var game Game

	s.gamesLock.RLock()
	game, has := s.games[id]
	s.gamesLock.RUnlock()

	if !has {
		return game, errors.New("game does not exist")
	}

	return game, nil
}

func (s *GameStorage) StartNew() uint64 {
	id := atomic.AddUint64(s.lastIndex, 1)

	game := Game{
		Id: id,
	}

	s.gamesLock.Lock()
	s.games[id] = game
	s.gamesLock.Unlock()

	return id
}