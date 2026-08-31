test:
    go test ./... -v -race -coverprofile=coverage.out

lint:
    golangci-lint run --fix
