# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git (required for go mod download)
RUN apk add --no-cache git ca-certificates

# Copy all source code
COPY . .

# Clean and regenerate dependencies
RUN rm -f go.sum && go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/.env.example .env

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]
