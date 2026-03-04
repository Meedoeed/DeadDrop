package storage

import "deaddrop/internal/models"

type Storage interface {
	Save(secret *models.Secret) error
	Get(id string) (*models.Secret, error)
	Delete(id string) error
	Close() error
}
