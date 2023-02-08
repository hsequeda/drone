include .env
export

.PHONY: run
run_server:
	go run ./cmd/server


.PHONY: test
test:
	go test ./... -v --race
