build:
	@go build -o bin/shop cmd/main.go
run: build
	@./bin/shop
create_db:
	@migrate create -ext sql -dir cmd/migrations -seq ${db_name}
migrate-up:
	@export $$(grep -v '^#' .env | xargs) && migrate -database $${POSTGRESQL_URL} -path cmd/migrations up
migrate-down:
	@export $$(grep -v '^#' .env | xargs) && migrate -database $${POSTGRESQL_URL} -path cmd/migrations down 1