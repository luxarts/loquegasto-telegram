package repository

import (
	"loquegasto-telegram/internal/defines"
	"sync"
)

type SessionsRepository interface {
	Add(userID int64, token string)
	Get(userID int64) (*string, error)
}

type sessionsRepository struct {
	values      map[int64]string
	valuesMutex sync.RWMutex
}

func NewSessionsRepository() SessionsRepository {
	return &sessionsRepository{
		values: make(map[int64]string, 0),
	}
}

func (s *sessionsRepository) Add(userID int64, token string) {
	s.valuesMutex.Lock()
	s.values[userID] = token
	s.valuesMutex.Unlock()
}

func (s *sessionsRepository) Get(userID int64) (*string, error) {
	token, ok := s.values[userID]

	if !ok {
		return nil, defines.ErrSessionNotFound
	}

	return &token, nil
}
