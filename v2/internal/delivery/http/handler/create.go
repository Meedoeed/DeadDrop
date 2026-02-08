package handler

import (
	"deaddrop/internal/usecase"
	"io"
	"log"
	"net/http"
	"path/filepath"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "File is too big or invalid form", http.StatusBadRequest)
		return
	}
	message := r.FormValue("message")
	ttl := r.FormValue("ttl")

	var fileData []byte
	var fileName string
	var fileExt string

	if file, fileHeader, err := r.FormFile("file"); err == nil {
		defer file.Close()

		if fileHeader.Size > 20<<20 {
			http.Error(w, "file is too big", http.StatusBadRequest)
			return
		}

		fileData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Cannot read file", http.StatusInternalServerError)
			return
		}

		fileName = filepath.Base(fileHeader.Filename)
		fileExt = filepath.Ext(fileHeader.Filename)
		mime := http.DetectContentType(fileData)
		allowedExts := map[string]bool{
			"image/jpeg":      true,
			"image/png":       true,
			"image/gif":       true,
			"application/pdf": true,
			"text/plain":      true,
		}

		if !allowedExts[mime] {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}
	}
	id, err := usecase.GenerateID(10)
	if err != nil {
		http.Error(w, "Cannot generate ID", http.StatusInternalServerError)
		return
	}
	log.Printf(
		"[INFO] POST /create | id=%s file=%s ttl=%s message=%s fileext=%s",
		id,
		fileName,
		ttl,
		message,
		fileExt,
	)
	http.Redirect(w, r, "/", 303)
}
