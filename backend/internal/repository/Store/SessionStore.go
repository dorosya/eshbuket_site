package store

import (
	"eshbuket/internal/Domain/models"
	"sync"
	"time"
)

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]models.Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{sessions: make(map[string]models.Session)}
}

//Получение сессии из памяти программы. Если просрочено = вернуть пустую сессию и False
func (s *SessionStore) Get(id string) (models.Session, bool) {
	s.mu.RLock()
	session, ok := s.sessions[id]
	s.mu.RUnlock()

	if !ok || session.Expires.Before(time.Now()) {
		s.Delete(id)
		return models.Session{}, false
	}
	return session, true
}

func (s *SessionStore) Set(id string, session models.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[id] = session
}

func (s *SessionStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
}
