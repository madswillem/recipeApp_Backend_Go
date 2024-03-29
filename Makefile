run_db:
	sudo systemctl start postgresql

run:
	go run cmd/main.go

test:
	go test -v -cover ./test/*_test.go

fmt:
	go fmt ./...

tidy:
	go mod tidy