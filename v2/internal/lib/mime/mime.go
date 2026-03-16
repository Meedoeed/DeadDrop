package mime

import (
	"net/http"
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"application/pdf": true,
	"text/plain":      true,
}

func IsMimeTypeAllowed(mimeType string) bool {
	return allowedMimeTypes[mimeType]
}

func DetectFromData(data []byte) string {
	if len(data) == 0 {
		return "application/octet-stream"
	}
	return http.DetectContentType(data)
}
