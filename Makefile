http_run:
	go run ./cmd/http/main.go

lint:
	golangci-lint run -v ./...

mock:

generate:
	go generate ./...
