FROM golang:1.19-alpine As builder
RUN apk --no-cache add ca-certificates
RUN mkdir /app_dir
COPY . /app_dir
WORKDIR /app_dir
RUN CGO_ENABLED=0 go build -o /app ./cmd/log_register/

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/log_register /log_register
CMD ["./log_register"]
