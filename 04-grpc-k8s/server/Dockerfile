FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY server/ ./server/
COPY proto/ ./proto/

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY certs/ ./certs/

CMD ["./server"]
