# ============================================
# BUILD STAGE - Compile Go application
# ============================================
FROM golang:1.25.6-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy and download dependencies first (better caching)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o /app/order-service \
    .

# ============================================
# FINAL STAGE - Minimal runtime image
# ============================================
FROM alpine:latest

LABEL maintainer="Order Tracking System" \
      description="Order Service - gRPC microservice for order management" \
      version="1.0"

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata && \
    # Create non-root user for security
    addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser && \
    # Create app directory
    mkdir -p /app && \
    chown -R appuser:appgroup /app

WORKDIR /app

# Copy the binary from builder
COPY --from=builder --chown=appuser:appgroup /app/order-service .

# Switch to non-root user
USER appuser

# Expose gRPC port
EXPOSE 50051

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:50051 || exit 1

# Run the application
ENTRYPOINT ["./order-service"]
