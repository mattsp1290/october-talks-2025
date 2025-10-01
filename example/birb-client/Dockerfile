# Build stage
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

RUN go install github.com/DataDog/orchestrion@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux orchestrion go build -o birb-client ./cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/birb-client .

# Run as non-root user
RUN adduser -D -u 1000 appuser
USER appuser

# Expose no ports (client doesn't listen)

ENTRYPOINT ["./birb-client"]