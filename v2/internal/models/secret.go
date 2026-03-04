package models

import (
	"time"
)

type Secret struct {
	ID        string
	Message   string
	FileData  []byte
	FileName  string
	FileExt   string
	Password  string
	ExpiresAt time.Time
}
