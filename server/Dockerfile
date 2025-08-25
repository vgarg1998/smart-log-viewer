# ===============================
# Builder Stage
# ===============================
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build for Linux AMD64 for portability
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-w -s" -o server ./cmd/server

# ===============================
# Production Stage
# ===============================
FROM alpine:latest AS production

# Runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /app

# Non-root user
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup

# Copy binary from builder
COPY --from=builder /app/server .

RUN chmod +x server && chown appuser:appgroup server

USER appuser

EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD curl -f http://localhost:8080/ || exit 1

CMD ["./server"]

# ===============================
# Development Stage (optional)
# ===============================
FROM golang:1.24-alpine AS development

RUN apk add --no-cache git curl ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Install air for hot reload
RUN go install github.com/cosmtrek/air@latest

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
