.PHONY: build run clean test docker-up docker-down


# ? Build all back services
build:
	@echo "Building services..."
	@go build -o bin/api-gateway ./services/api-gateway
	@go build -o bin/auth-service ./services/auth-service
	@go build -o bin/users-service ./services/users-service
	@go build -o bin/teams-service ./services/teams-service
	@go build -o bin/boards-service ./services/boards-service
	@go build -o bin/tasks-service ./services/tasks-service
	@echo "Build complete!"


# ? Run back services locally
run-gateway:
	@go run services/api-gateway/*.go

run-auth:
	@go run services/auth-service/main.go

run-users:
	@go run services/users-service/main.go

run-teams:
	@go run services/teams-service/main.go

run-boards:
	@go run services/boards-service/main.go

run-tasks:
	@go run services/tasks-service/main.go


# ? Clean build artifacts
clean:
	@rm -rf bin/
	@echo "Cleaned!"


# ? Docker commands
docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down

docker-logs:
	@docker-compose logs -f


# ? Install dependencies
deps:
	@go mod download
	@go mod tidy


# ? Format code
fmt:
	@go fmt ./...


# ? Run tests
test:
	@go test ./...


run-tests:
	@powershell -NoProfile -ExecutionPolicy Bypass -File tests\run-tests.ps1 -RunAll


# ? Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@swag init -g services/api-gateway/main.go -o services/api-gateway/docs
	@echo "Swagger documentation generated in services/api-gateway/docs/"
