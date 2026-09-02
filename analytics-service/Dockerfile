# ============================================
# BUILD STAGE - Compile Go application
# ============================================
FROM golang:1.25.6-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy analytics-service module files for dependency caching
COPY analytics-service/go.mod analytics-service/go.sum analytics-service/

# Create a minimal go.work that only references analytics-service
RUN printf 'go 1.25.6\n\nuse (\n\t./analytics-service\n)\n' > go.work

# Download dependencies
RUN cd analytics-service && go mod download

# Copy full analytics-service source
COPY analytics-service/ analytics-service/

# Build the application from workspace root
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o /app/analytics-service \
    ./analytics-service/.

# ============================================
# FINAL STAGE - Minimal runtime image
# ============================================
FROM alpine:latest

LABEL maintainer="Order Tracking System" \
      description="Analytics Service - HTTP REST API for search and analytics" \
      version="1.0"

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata curl && \
    # Create non-root user for security
    addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser && \
    # Create app directory
    mkdir -p /app && \
    chown -R appuser:appgroup /app

WORKDIR /app

# Copy the binary from builder
COPY --from=builder --chown=appuser:appgroup /app/analytics-service .

# Switch to non-root user
USER appuser

# Expose HTTP port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
ENTRYPOINT ["./analytics-service"]
