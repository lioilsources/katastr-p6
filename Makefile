.PHONY: up down backend-run backend-test backend-build app-run app-analyze app-test

# Docker
up:
	docker compose up -d

down:
	docker compose down

# Backend (Go)
backend-run:
	cd backend && go run cmd/server/main.go

backend-test:
	cd backend && go test -v ./...

backend-build:
	cd backend && go build -o bin/server cmd/server/main.go

# App (Flutter)
app-run:
	cd app && flutter run

app-analyze:
	cd app && flutter analyze

app-test:
	cd app && flutter test
