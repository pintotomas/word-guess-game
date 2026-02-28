package store

import (
	"sync"

	"github.com/fleetdm/wordgame/game"
	"github.com/pkg/errors"
)

var ErrGameNotFound = errors.New("game not found")

type Store struct {
	mu    sync.RWMutex
	games map[string]*game.Game
}

func New() *Store {
	return &Store{
		games: make(map[string]*game.Game),
	}
}

func (s *Store) Save(id string, g *game.Game) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.games[id] = g
}

func (s *Store) Get(id string) (*game.Game, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	g, ok := s.games[id]
	if !ok {
		return nil, ErrGameNotFound
	}
	return g, nil
}

func (s *Store) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.games, id)
}
