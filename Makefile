run:
	docker-compose up -d --build
stop:
	docker-compose down

# Backend commands
backend-run:
	cd backend && go run cmd/http/main.go

backend-build:
	cd backend && go build -o bin/main cmd/http/main.go

backend-test:
	cd backend && go test ./...

backend-deps:
	cd backend && go mod tidy

# Database commands
db-start:
	docker-compose up -d db

db-stop:
	docker-compose stop db

# Frontend commands
frontend-run:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

frontend-deps:
	cd frontend && npm install

# Clean commands
clean:
	docker-compose down -v
	rm -rf backend/bin

.PHONY: run stop backend-run backend-build backend-test backend-deps db-start db-stop frontend-run frontend-build frontend-deps clean