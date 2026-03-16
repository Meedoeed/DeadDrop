package handler

import (
	"deaddrop/internal/assets"
	"deaddrop/internal/usecase"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type CreateHandler struct {
	secretUseCase *usecase.SecretUseCase
}

// NewCreateHandler - конструктор для нашего хэндлера. Мы внедрили зависимость от uc
func NewCreateHandler(uc *usecase.SecretUseCase) *CreateHandler {
	return &CreateHandler{
		secretUseCase: uc,
	}
}

func (h *CreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Здесь проверяем метод запроса - это слой доставки
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Ограничение размера тела запроса для защиты от DoS-атак. Мы обсуждали это выше
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	// Парсим multipart форму и ограничиваем размер
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "File is too big or invalid form", http.StatusBadRequest)
		return
	}

	// Достаем данные из формы, готовимся к передаче в uc
	message := r.FormValue("message")
	ttl := r.FormValue("ttl")

	var fileData []byte
	var fileName string
	var fileExt string

	// Если файл в наличии пытаемся его обработать
	if file, fileHeader, err := r.FormFile("file"); err == nil {
		defer file.Close()

		// Прооверяем размер файла через fileHeader.Size
		if fileHeader.Size > 20<<20 {
			http.Error(w, "File is too big", http.StatusBadRequest)
			return
		}

		// Читаем содержимое файла
		fileData, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Cannot read file", http.StatusInternalServerError)
			return
		}

		// Проверка чтобы имя файла не включало в себя пути, оно не должно иметь вид "/home/.." или "./..." - это опасно
		fileName = filepath.Base(fileHeader.Filename)
		fileExt = filepath.Ext(fileHeader.Filename)
	}

	// Готовим структуру для отправки в uc
	createReq := &usecase.CreateSecretRequest{
		Message:  message,
		TTL:      ttl,
		FileData: fileData,
		FileName: fileName,
		FileExt:  fileExt,
	}

	// Отправляем данные в uc-слой
	resp, err := h.secretUseCase.Create(createReq)
	if err != nil {
		log.Printf("[ERROR] Failed to create secret: %v", err)
		http.Error(w, "Failed to create secret: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Шаг 7: Логирование успешного создания
	log.Printf("[INFO] POST /create | id=%s ttl=%s message='%s' file=%s",
		resp.ID, ttl, message, fileName)

	// Перенаправляем пользователя на страницу с результатом
	// Позже здесь будет редирект на /secret/{id}, а пароль нужно будет показать на отдельной странице.
	// Пока же мы просто вернём ID и пароль в теле ответа для наглядности.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.WriteHeader(http.StatusCreated)
	// Подготавливаем данные для шаблона
	data := struct {
		Title     string
		ID        string
		Password  string
		ExpiresAt string
	}{
		Title:     "DeadDrop — ячейка создана",
		ID:        resp.ID,
		Password:  resp.Password,
		ExpiresAt: resp.ExpiresAt.Format(time.RFC1123),
	}

	// Рендерим шаблон
	err = assets.ResultTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("[ERROR] Failed to render template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
