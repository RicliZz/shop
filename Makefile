build:
	@go build -o bin/shop cmd/main.go
run: build
	@go run cmd/main.go