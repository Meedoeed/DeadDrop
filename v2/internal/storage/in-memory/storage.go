package inmemory

import (
	"deaddrop/internal/models"
	"sync"
)

type Storage struct {
	mu      sync.RWMutex
	secrets map[string]*models.Secret
}

func NewStorage() *Storage {
	return &Storage{
		secrets: make(map[string]*models.Secret),
	}
}

func (s *Storage) Save(secret *models.Secret) error {
	s.mu.Lock()         // Блокируем доступ для всех остальных горутин
	defer s.mu.Unlock() // Гарантированно разблокируем при выходе из функции

	s.secrets[secret.ID] = secret
	return nil
}

func (s *Storage) Get(id string) (*models.Secret, error) {
	s.mu.RLock()         // Блокируем только запись, чтение разрешено другим
	defer s.mu.RUnlock() // Снимаем блокировку чтения

	secret, ok := s.secrets[id]
	if !ok {
		return nil, nil // Секрет не найден
	}
	return secret, nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.secrets, id)
	return nil
}

func (s *Storage) Close() error {
	return nil
}
