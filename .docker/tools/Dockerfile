FROM golang:1.19-alpine As builder
RUN apk --no-cache add ca-certificates
RUN mkdir /app_dir
COPY . /app_dir
WORKDIR /app_dir
