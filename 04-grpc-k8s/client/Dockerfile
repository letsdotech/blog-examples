# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY client/ ./client/
COPY proto/ ./proto/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o client ./client

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/client .
COPY certs/ ./certs/

# Run the binary
CMD ["./client"]


#letsdotech/grpc-client:latest
