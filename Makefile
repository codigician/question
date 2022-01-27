up:
	docker build --progress=plain -t question-api .
	docker run question-api

run:
	go run .

unit-test:
	go test ./... -v -short

test:
	go test ./... -v

code-coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out | grep total | awk '{print $3}'

lint:
	golangci-lint run

mockgen:
	mockgen -destination=mocks/mock_service.go -package mocks -source=handler.go
	mockgen -destination=mocks/mock_repository.go -package mocks -source=service.go