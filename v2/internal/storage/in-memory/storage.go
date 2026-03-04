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
