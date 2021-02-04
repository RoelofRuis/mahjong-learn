package main

import (
	"errors"
	"github.com/roelofruis/mahjong-learn/mahjong"
	"sync"
	"sync/atomic"
)

func NewGameStorage() *GameStorage {
	return &GameStorage{
		gamesLock: sync.RWMutex{},
		games:     make(map[uint64]*mahjong.Game),
		lastIndex: new(uint64),
	}
}

type GameStorage struct {
	gamesLock sync.RWMutex
	games     map[uint64]*mahjong.Game

	lastIndex *uint64
}

func (s *GameStorage) Get(id uint64) (*mahjong.Game, error) {
	var g *mahjong.Game

	s.gamesLock.RLock()
	g, has := s.games[id]
	s.gamesLock.RUnlock()

	if !has {
		return nil, errors.New("game does not exist")
	}

	return g, nil
}

func (s *GameStorage) StartNew() (uint64, error) {
	id := atomic.AddUint64(s.lastIndex, 1)

	m, err := mahjong.NewGame(id)
	if err != nil {
		return id, err
	}

	s.gamesLock.Lock()
	s.games[id] = m
	s.gamesLock.Unlock()

	return id, nil
}
