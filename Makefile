include .env
export

.PHONY: run
run_server:
	go run .


.PHONY: test
test:
	go test ./... -v --race
