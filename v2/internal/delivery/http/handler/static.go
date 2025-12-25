package handler

import (
	"deaddrop/internal/assets"
	"net/http"
)

var StaticHandler http.Handler = http.FileServer(
	http.FS(assets.StaticFS),
)
