# Use lightweight Go base image for building
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary file (static linking)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o opencode-ssh-mcp .

# Use lightweight alpine as final image
FROM alpine:latest

# Install necessary dependencies
RUN apk --no-cache add ca-certificates openssh-client

# Create non-root user
RUN adduser -D -s /bin/sh opencode

# Create necessary directories
RUN mkdir -p /home/opencode/.ssh && \
    chown -R opencode:opencode /home/opencode

# Copy binary from builder
COPY --from=builder /app/opencode-ssh-mcp /usr/local/bin/opencode-ssh-mcp

# Set user permissions
USER opencode

# Set entrypoint
ENTRYPOINT ["opencode-ssh-mcp"]