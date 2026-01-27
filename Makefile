http_run:
	go run ./cmd/api/main.go

lint:
	golangci-lint run -v ./...

mock:

generate:
	go generate ./...
