package usecase

import (
	"errors"
	"fmt"
	"time"

	"deaddrop/internal/lib/generator"
	"deaddrop/internal/lib/hash"
	"deaddrop/internal/lib/mime"
	"deaddrop/internal/models"
	"deaddrop/internal/storage"
)

type SecretUseCase struct {
	storage storage.Storage
}

type CreateSecretRequest struct {
	Message  string
	TTL      string
	FileData []byte
	FileName string
	FileExt  string
}

type CreateSecretResponse struct {
	ID        string
	Password  string
	ExpiresAt time.Time
}

type GetSecretRequest struct {
	ID       string
	Password string
}

type GetSecretResponse struct {
	ID        string
	Message   string
	FileName  string
	FileExt   string
	HasFile   bool
	ExpiresAt time.Time
}

type GetFileRequest struct {
	ID string
}

type GetFileResponse struct {
	FileData []byte
	FileName string
	FileExt  string
	Message  string
}

func NewSecretUseCase(s storage.Storage) *SecretUseCase {
	return &SecretUseCase{
		storage: s,
	}
}

func (uc *SecretUseCase) Create(req *CreateSecretRequest) (*CreateSecretResponse, error) {
	id, err := generator.GenerateID(10)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ID: %w", err)
	}

	rawPassword, err := generator.GeneratePassword(12)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	hashedPassword, err := hash.HashPassword(rawPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	ttlSeconds := 3600 // значение по умолчанию - 1 час
	if req.TTL != "" {
		var parsedTTL int
		_, err := fmt.Sscanf(req.TTL, "%d", &parsedTTL)
		if err == nil && parsedTTL > 0 {
			ttlSeconds = parsedTTL
		}
	}

	expiresAt := time.Now().Add(time.Duration(ttlSeconds) * time.Hour)

	if len(req.FileData) > 0 {
		detectedMime := mime.DetectFromData(req.FileData)
		if !mime.IsMimeTypeAllowed(detectedMime) {
			return nil, errors.New("file type not allowed")
		}
	}

	secret := &models.Secret{
		ID:        id,
		Message:   req.Message,
		FileData:  req.FileData,
		FileName:  req.FileName,
		FileExt:   req.FileExt,
		Password:  hashedPassword,
		ExpiresAt: expiresAt,
	}

	if err := uc.storage.Save(secret); err != nil {
		return nil, fmt.Errorf("failed to save secret: %w", err)
	}

	return &CreateSecretResponse{
		ID:        id,
		Password:  rawPassword,
		ExpiresAt: expiresAt,
	}, nil
}

func (uc *SecretUseCase) GetSecret(req *GetSecretRequest) (*GetSecretResponse, error) {
	secret, err := uc.storage.Get(req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}
	if secret == nil {
		return nil, errors.New("secret not found")
	}

	if time.Now().After(secret.ExpiresAt) {
		uc.storage.Delete(req.ID)
		return nil, errors.New("secret has expired")
	}

	if !hash.CheckPasswordHash(req.Password, secret.Password) {
		return nil, errors.New("invalid password")
	}

	return &GetSecretResponse{
		ID:        secret.ID,
		Message:   secret.Message,
		FileName:  secret.FileName,
		FileExt:   secret.FileExt,
		HasFile:   len(secret.FileData) > 0,
		ExpiresAt: secret.ExpiresAt,
	}, nil
}

func (uc *SecretUseCase) GetFile(req *GetFileRequest) (*GetFileResponse, error) {
	secret, err := uc.storage.Get(req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}
	if secret == nil {
		return nil, errors.New("secret not found")
	}

	// Проверяем, не истёк ли срок
	if time.Now().After(secret.ExpiresAt) {
		uc.storage.Delete(req.ID)
		return nil, errors.New("secret has expired")
	}

	// Проверяем, есть ли файл
	if len(secret.FileData) == 0 {
		return nil, errors.New("no file attached to this secret")
	}

	return &GetFileResponse{
		FileData: secret.FileData,
		FileName: secret.FileName,
		FileExt:  secret.FileExt,
		Message:  secret.Message,
	}, nil
}
