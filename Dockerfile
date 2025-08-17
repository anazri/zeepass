# Build stage
FROM golang:1.24.2-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies)
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o zeepass \
    cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and tzdata for timezone
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S zeepass && \
    adduser -u 1001 -S zeepass -G zeepass

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/zeepass .

# Copy static assets and templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Change ownership to non-root user
RUN chown -R zeepass:zeepass /app

# Switch to non-root user
USER zeepass

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8080

# Run the binary
CMD ["./zeepass"]