package handler

import (
	"deaddrop/internal/assets"
	"deaddrop/internal/lib/auth"
	"deaddrop/internal/lib/url"
	"deaddrop/internal/usecase"
	"log"
	"net/http"
	"time"
)

type SecretHandler struct {
	secretUseCase *usecase.SecretUseCase
}

func NewSecretHandler(uc *usecase.SecretUseCase) *SecretHandler {
	return &SecretHandler{
		secretUseCase: uc,
	}
}

func (h *SecretHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, action, ok := url.ParseSecretPath(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}

	switch action {
	case "view":
		switch r.Method {
		case http.MethodGet:
			h.showPasswordForm(w, r, id)
		case http.MethodPost:
			h.handleSecretAccess(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "download":
		// Это запрос на скачивание файла — только GET
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleFileDownload(w, r, id)
	default:
		http.NotFound(w, r)
	}
}

func (h *SecretHandler) showPasswordForm(w http.ResponseWriter, r *http.Request, id string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := struct {
		ID    string
		Title string
	}{
		ID:    id,
		Title: "Введите пароль",
	}

	if err := assets.SecretPassTemplate.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("[ERROR] Failed to render password template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *SecretHandler) handleSecretAccess(w http.ResponseWriter, r *http.Request, id string) {
	password := r.FormValue("password")
	if password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	req := &usecase.GetSecretRequest{
		ID:       id,
		Password: password,
	}

	resp, err := h.secretUseCase.GetSecret(req)
	if err != nil {
		log.Printf("[ERROR] Failed to get secret %s: %v", id, err)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := struct {
			ID    string
			Title string
			Error string
		}{
			ID:    id,
			Title: "Введите пароль",
			Error: err.Error(),
		}
		assets.SecretPassTemplate.ExecuteTemplate(w, "layout", data)
		return
	}

	auth.SetAuthCookie(w, id)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	contentData := struct {
		ID        string
		Message   string
		FileName  string
		HasFile   bool
		ExpiresAt string
		Title     string
	}{
		ID:        resp.ID,
		Message:   resp.Message,
		FileName:  resp.FileName,
		HasFile:   resp.HasFile,
		ExpiresAt: resp.ExpiresAt.Format(time.RFC1123),
		Title:     "DeadDrop",
	}

	if err := assets.SecretContentTemplate.ExecuteTemplate(w, "layout", contentData); err != nil {
		log.Printf("[ERROR] Failed to render content template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *SecretHandler) handleFileDownload(w http.ResponseWriter, r *http.Request, id string) {
	// 1: Проверяем авторизацию через подписанную куку
	if !auth.CheckAuthCookie(r, id) {
		log.Printf("[WARN] Unauthorized download attempt for secret %s from %s", id, r.RemoteAddr)
		http.Redirect(w, r, "/secret/"+id, http.StatusSeeOther)
		return
	}

	// 2: Создаём запрос для Use Case
	req := &usecase.GetFileRequest{
		ID: id,
	}

	// 3: Вызываем Use Case для получения данных файла
	resp, err := h.secretUseCase.GetFile(req)
	if err != nil {
		log.Printf("[ERROR] Failed to get file for secret %s: %v", id, err)

		switch {
		case err.Error() == "secret not found" || err.Error() == "secret has expired":
			http.NotFound(w, r)
		case err.Error() == "no file attached to this secret":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			data := struct {
				ID    string
				Title string
				Error string
			}{
				ID:    id,
				Title: "Содержимое ячейки",
				Error: "К этой ячейке не прикреплён файл",
			}
			assets.SecretContentTemplate.ExecuteTemplate(w, "layout", data)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	fileName := resp.FileName
	if fileName == "" {
		fileName = "secret_" + id + resp.FileExt
	}

	// Content-Disposition: attachment заставляет браузер скачать файл, а не открывать его
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")

	// Определяем Content-Type по расширению или данным файла
	mimeType := http.DetectContentType(resp.FileData)
	w.Header().Set("Content-Type", mimeType)

	// 5: Отдаём файл
	w.WriteHeader(http.StatusOK)
	w.Write(resp.FileData)

	log.Printf("[INFO] File %s downloaded for secret %s", fileName, id)
}
