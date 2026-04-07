package in_mem

import (
	"media-equipment-tracker/internal/config"
	"sync"
	"time"
)

var tokenRepository TokenRepository
var tokenRepositoryMutex sync.Mutex

type TokenRepository interface {
	Check(token string) bool
	Add(token string) bool
	Delete(token string)
	cleanupExpired()
}

type tokenEntry struct {
	expiresAt time.Time
}

type mapTokenRepositoryWithTTL struct {
	mu       sync.RWMutex
	tokens   map[string]tokenEntry
	ttl      time.Duration
	stopChan chan struct{}
}

func GetTokenRepository() TokenRepository {
	tokenRepositoryMutex.Lock()
	defer tokenRepositoryMutex.Unlock()

	if tokenRepository != nil {
		return tokenRepository
	}

	tokenRepository = &mapTokenRepositoryWithTTL{
		tokens:   make(map[string]tokenEntry),
		ttl:      config.AccessTokenDuration,
		stopChan: make(chan struct{}),
	}
	// горутина для очистки просроченных токенов
	go tokenRepository.cleanupExpired()
	return tokenRepository
}

func (r *mapTokenRepositoryWithTTL) Check(token string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.tokens[token]
	if !exists {
		return false
	}

	if time.Now().After(entry.expiresAt) {
		return false
	}
	return true
}

func (r *mapTokenRepositoryWithTTL) Add(token string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if entry, exists := r.tokens[token]; exists {
		if time.Now().Before(entry.expiresAt) {
			return false // активный токен уже существует
		}
	}

	r.tokens[token] = tokenEntry{
		expiresAt: time.Now().Add(r.ttl),
	}
	return true
}

func (r *mapTokenRepositoryWithTTL) Delete(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tokens, token)
}

// периодически удаляет просроченные токены
func (r *mapTokenRepositoryWithTTL) cleanupExpired() {
	ticker := time.NewTicker(r.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.removeExpired()
		case <-r.stopChan:
			return
		}
	}
}

func (r *mapTokenRepositoryWithTTL) removeExpired() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for token, entry := range r.tokens {
		if now.After(entry.expiresAt) {
			delete(r.tokens, token)
		}
	}
}

func (r *mapTokenRepositoryWithTTL) Stop() {
	close(r.stopChan)
}
