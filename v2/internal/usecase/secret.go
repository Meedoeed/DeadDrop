package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateID(lenght int) (string, error) {
	bytes := make([]byte, lenght)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("error in ID generation: %s", err)
	}
	id := base64.RawURLEncoding.EncodeToString(bytes)
	return id, nil
}
