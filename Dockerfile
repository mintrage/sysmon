FROM golang:1.26.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o sysmon-server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/sysmon-server .

EXPOSE 8080

CMD ["./sysmon-server"]