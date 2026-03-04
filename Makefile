.PHONY: fmt vet lint test run build

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

test:
	go test ./... -count=1

run:
	go run cmd/app/main.go

build:
	go build -o bin/app cmd/app/main.go

pre-commit:
	pre-commit run --all-files