.PHONY: test coverage

test:
	go test -v ./...

coverage:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o coverage.html

lint:
	golangci-lint run ./...
