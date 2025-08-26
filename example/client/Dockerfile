# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# First, copy the local SDK dependency if using replace directive
# This assumes you'll copy the SDK to the build context
COPY ag-ui-go-sdk /ag-ui-go-sdk

# Copy all source code including go.mod
COPY . .

# Update replace directive to use container path
RUN sed -i 's|/Users/punk1290/git/ag-ui/go-sdk|/ag-ui-go-sdk|' go.mod

# Download dependencies
RUN go mod download

# Build the headless application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o headless cmd/headless/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/headless .

ENV SERVER_HOST="host.docker.internal"

# Set the entrypoint
ENTRYPOINT ["./headless"]

# Default CMD (can be overridden)
CMD ["endpoint"]