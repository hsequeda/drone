FROM golang:1.19-alpine As builder
RUN apk --no-cache add ca-certificates
RUN mkdir /app_dir
COPY . /app_dir
WORKDIR /app_dir
RUN apk add build-base
RUN go mod tidy
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.1
