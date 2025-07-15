FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN go build -o server ./cmd/server.go

FROM alpine:3.22.0

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/internal/storages/postgresql/migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["./server"]