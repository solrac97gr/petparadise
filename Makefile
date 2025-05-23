# Run commands
run:
	docker-compose up -d --build
stop:
	docker-compose down

# Test commands
integration:
	docker-compose -f backend/test/docker-compose.test.yml up -d --build
	sleep 10
	cd backend && TEST_DATABASE_URL=postgres://testuser:testpass@localhost:5433/petparadise_test?sslmode=disable TEST_API_URL=http://localhost:3001/api go test ./test/integration/... -v
	docker-compose -f backend/test/docker-compose.test.yml down

# Clean commands
clean:
	docker-compose down -v
	rm -rf backend/bin

.PHONY: run stop integration clean