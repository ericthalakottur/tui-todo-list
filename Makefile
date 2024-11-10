.PHONY: sqlc
sqlc: sqlc.yml query.sql schema.sql
ifeq ($(shell which sqlc), '')
	$(error "sqlc not installed")
endif
	sqlc generate

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build: tidy sqlc
	go build -o build/app .

.PHONY: run
run: tidy sqlc
	go run .
