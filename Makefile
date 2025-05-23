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

# Test commands
integration-test:
	cd backend && go test ./test/integration/... -v

integration-test-docker:
	docker-compose -f backend/test/docker-compose.test.yml up -d --build
	sleep 10
	cd backend && TEST_DATABASE_URL=postgres://testuser:testpass@localhost:5433/petparadise_test?sslmode=disable TEST_API_URL=http://localhost:3001/api/users go test ./test/integration/... -v
	docker-compose -f backend/test/docker-compose.test.yml down

integration-test-setup:
	docker-compose -f backend/test/docker-compose.test.yml up -d --build

integration-test-cleanup:
	docker-compose -f backend/test/docker-compose.test.yml down -v

# Clean commands
clean:
	docker-compose down -v
	rm -rf backend/bin

.PHONY: run stop backend-run backend-build backend-test backend-deps db-start db-stop frontend-run frontend-build frontend-deps clean