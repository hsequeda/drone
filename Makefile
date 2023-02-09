include .env
export

.PHONY: run_server
run_server:
	@docker-compose run --rm drone_server

.PHONY: tools
tools:
	@docker-compose build tools

.PHONY: run_log_register
run_log_register:
	@docker-compose run --rm drone_log_register

.PHONY: test
test:
	@docker-compose run --rm tools go test ./... -v --race


.PHONY: lint
lint:
	@docker-compose run --rm tools golangci-lint run
