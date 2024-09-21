run_db:
    sudo systemctl start postgresql

run_docker:
    sudo systemctl start docker

run:
	go run cmd/main/main.go

run_dev:
	go run cmd/dev/dev.go


run_api:
	go run cmd/api.go

run_web:
	go run cmd/web.go

test:
	go test -v -cover ./test/*_test.go

fmt:
	go fmt ./...

tidy:
	go mod tidy
