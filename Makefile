include .env
export

.PHONY: run_server
run_server:
	go run ./cmd/server

.PHONY: run_log_register
run_log_register:
	go run ./cmd/log_register

.PHONY: test
test:
	go test ./... -v --race
