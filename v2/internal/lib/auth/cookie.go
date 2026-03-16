package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
)

// ВНИМАНИЕ! Это временное решение для обучения.
// В реальном проекте секретный ключ НИКОГДА не должен храниться в коде!
// Он должен загружаться из переменных окружения или безопасного хранилища.
var (
	secretKey = []byte("A8zCb7GtY$!") // TODO: вынести в конфиг!
)

func createSignature(data string) string {
	h := hmac.New(sha256.New, secretKey)

	h.Write([]byte(data))

	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func validateSignature(data, signature string) bool {
	expectedSignature := createSignature(data)

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func SetAuthCookie(w http.ResponseWriter, secretID string) {
	signature := createSignature(secretID)

	value := secretID + ":" + signature

	cookie := &http.Cookie{
		Name:     "secret_auth_" + secretID,
		Value:    value,
		Path:     "/secret/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600,
	}

	http.SetCookie(w, cookie)
}

func CheckAuthCookie(r *http.Request, secretID string) bool {
	cookie, err := r.Cookie("secret_auth_" + secretID)
	if err != nil {
		return false
	}

	parts := strings.Split(cookie.Value, ":")
	if len(parts) != 2 {
		return false
	}

	cookieID := parts[0]
	signature := parts[1]

	if cookieID != secretID {
		return false
	}

	return validateSignature(secretID, signature)
}

func ClearAuthCookie(w http.ResponseWriter, secretID string) {

	cookie := &http.Cookie{
		Name:     "secret_auth_" + secretID,
		Value:    "",
		Path:     "/secret/",
		HttpOnly: true,
		MaxAge:   -1,
	}

	http.SetCookie(w, cookie)
}
