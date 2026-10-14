# Notez — orquestração dos apps (api, files, web)

-include .env
export

API_DIR   := apps/api
FILES_DIR := apps/files
WEB_DIR   := apps/web
BIN_DIR   := bin

.DEFAULT_GOAL := help
.PHONY: help setup install database down dev dev-api dev-files dev-web \
        build build-api build-files build-web test fmt migration clean

help:
	@echo "Notez — comandos disponíveis:"
	@echo ""
	@echo "  make setup        cria .env e apps/web/.env a partir dos exemplos"
	@echo "  make install      instala dependências (go mod + bun install)"
	@echo "  make database     sobe o Postgres via docker compose"
	@echo "  make down         derruba o Postgres"
	@echo "  make dev          roda api + files + web juntos"
	@echo "  make dev-api      roda apenas a API"
	@echo "  make dev-files    roda apenas o serviço de files"
	@echo "  make dev-web      roda apenas o front"
	@echo "  make build        compila api + files + front"
	@echo "  make test         roda os testes dos serviços Go"
	@echo "  make fmt          formata o código Go"
	@echo "  make migration name=nome   cria uma migration (goose)"
	@echo "  make clean        remove artefatos de build"

setup:
	@test -f .env || cp .env.example .env
	@test -f $(WEB_DIR)/.env || cp $(WEB_DIR)/.env.example $(WEB_DIR)/.env
	@echo "arquivos .env prontos — edite-os com suas credenciais"

install:
	cd $(API_DIR) && go mod download
	cd $(FILES_DIR) && go mod download
	cd $(WEB_DIR) && bun install

database:
	docker compose up -d database

down:
	docker compose down

dev:
	@trap 'kill 0' INT TERM EXIT; \
	( cd $(FILES_DIR) && go run ./cmd/main.go ) & \
	( cd $(API_DIR) && go run ./cmd/main.go ) & \
	( cd $(WEB_DIR) && bun run dev ) & \
	wait

dev-api:
	cd $(API_DIR) && go run ./cmd/main.go

dev-files:
	cd $(FILES_DIR) && go run ./cmd/main.go

dev-web:
	cd $(WEB_DIR) && bun run dev

build: build-api build-files build-web

build-api:
	mkdir -p $(BIN_DIR) && cd $(API_DIR) && go build -o ../../$(BIN_DIR)/api ./cmd/main.go

build-files:
	mkdir -p $(BIN_DIR) && cd $(FILES_DIR) && go build -o ../../$(BIN_DIR)/files ./cmd/main.go

build-web:
	cd $(WEB_DIR) && bun run build

test:
	cd $(API_DIR) && go test ./...
	cd $(FILES_DIR) && go test ./...

fmt:
	cd $(API_DIR) && gofmt -w .
	cd $(FILES_DIR) && gofmt -w .

migration:
	@test -n "$(name)" || ( echo "uso: make migration name=nome_da_migracao"; exit 1 )
	cd $(API_DIR) && goose -dir=./internal/migration create $(name) sql

clean:
	rm -rf $(BIN_DIR) $(WEB_DIR)/.svelte-kit $(WEB_DIR)/build
