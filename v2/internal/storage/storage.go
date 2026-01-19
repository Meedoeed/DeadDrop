package storage

import "time"

type storage struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	FileName  string    `json:"file_name"`
	FileData  []byte    `json:"-"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type CreateSecretRequest struct {
	Message  string
	FileData []byte
	FileName string
	Password string
	TTLHours int
}
