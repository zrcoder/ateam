.PHONY: build test lint fmt

run:
	go run .

test:
	go test ./...

fmt:
	go mod tidy
	go fmt ./...

lint: fmt
	go fix ./... && go vet ./...

check: fmt lint test