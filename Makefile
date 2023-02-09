include .env
export

.PHONY: run_server
run_server:
	go run ./cmd/server

.PHONY: tools
tools:
	@docker-compose build -q tools

.PHONY: run_log_register
run_log_register:
	go run ./cmd/log_register

.PHONY: test
test:
	@docker-compose run --rm tools "go test ./... -v --race"


.PHONY: lint
lint:
	@docker-compose run --rm tools "golangci-lint run"
