# Makefile

# Start PostgreSQL service
run_db:
	sudo systemctl start postgresql

# Start Docker service
run_docker:
	sudo systemctl start docker

# Run the main Go application
run:
	go run cmd/main/main.go

# Run the dev version of the Go application
run_dev:
	go run cmd/dev/dev.go

# Run the API version of the Go application
run_api:
	go run cmd/api.go

# Run the web version of the Go application
run_web:
	go run cmd/web.go

# Run Go tests with verbose output and coverage report
test:
	go test -v -cover ./test/*_test.go

# Format Go code
fmt:
	go fmt ./...

# Clean up Go module dependencies
tidy:
	go mod tidy
