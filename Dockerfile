FROM golang:1.23.2-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

FROM debian:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/internal/migrations /migrations

CMD ["./main"]
