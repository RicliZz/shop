build:
	@go build -o bin/shop cmd/main.go
run: build
	@./bin/shop