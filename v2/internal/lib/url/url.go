package url

import (
	"strings"
)

func ParseSecretPath(path string) (id, action string, ok bool) {
	const prefix = "/secret/"

	if !strings.HasPrefix(path, prefix) {
		return "", "", false
	}

	remaining := strings.TrimPrefix(path, prefix)
	parts := strings.Split(remaining, "/")

	if len(parts) == 0 || parts[0] == "" {
		return "", "", false
	}

	id = parts[0]

	if len(parts) == 1 {
		return id, "view", true
	}

	if len(parts) == 2 && parts[1] == "download" {
		return id, "download", true
	}

	return "", "", false
}
