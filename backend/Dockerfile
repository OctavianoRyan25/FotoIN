# Multi-stage Dockerfile for FotoIN Backend API
FROM golang:alpine AS builder

WORKDIR /app

# Install git and certificates
RUN apk add --no-cache git ca-certificates tzdata

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/api ./cmd/api

# Final runtime image
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/api /app/api

EXPOSE 8080

CMD ["/app/api"]
