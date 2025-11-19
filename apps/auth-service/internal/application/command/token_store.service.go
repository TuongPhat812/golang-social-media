package command

import "sync"

type TokenStore struct {
	mu     sync.RWMutex
	tokens map[string]string
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]string),
	}
}

func (s *TokenStore) Save(token, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = userID
}

func (s *TokenStore) GetUserID(token string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.tokens[token]
	return id, ok
}
