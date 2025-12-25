APP_NAME=deaddrop
CMD_DIR=cmd/api

.PHONY: run build clean

## Запуск сервера (гарантирует cwd = корень проекта)
run:
	cd $(CURDIR) && go run ./cmd/api

## Сборка бинарника
build:
	go build -o bin/$(APP_NAME) ./$(CMD_DIR)

## Очистка артефактов
clean:
	rm -rf bin
