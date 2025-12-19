package handler

import (
	templates "deaddrop/internal/delivery/http/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := struct {
		Message string
		Title   string
	}{
		Message: "DeadDrop is Alive",
		Title:   "DeadDrop",
	}

	err := templates.HomeTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Render error", http.StatusInternalServerError)
	}
}
