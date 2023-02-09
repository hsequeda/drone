# Drone

## Build

Execute `make build` to build the base images for `test` and `run`.

## Test

Run `make test` for execute the `e2e`, `unit` and `integration` test.
Run `make lint` to execute the `go`

## Run

### API Server

#### Configuration

A configuration example can be found in `env.dist`.

`HTTP_SERVER_ADDR`: HTTP server address.
`UPLOAD_SIZE`: Max upload size for Medications (In Mb).

#### Setup

Rename env.dist to .env (the configuration can be modified).
Define the `port mapping` and the `shared volumes` in `docker-compose`.

#### Execute

Run `make run_server`.

### Log Register

#### Configuration
A configuration example can be found in `env.dist`.

`LOG_REGISTER_INTERVAL`: Amount of time (in seconds) that the `log_register` require to re-run the job.

#### Setup

Rename env.dist to .env (the configuration can be modified).
Define the `shared volumes` in `docker-compose`.

#### Execute

Run `make log_register`.
