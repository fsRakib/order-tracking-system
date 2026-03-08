# ============================================
# BUILD STAGE - Compile Go application
# ============================================
FROM golang:1.25.6-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy the shared pb module
COPY pb/go.mod pb/go.sum pb/
COPY pb/stock/ pb/stock/

# Copy the stock-service module
COPY stock-service/go.mod stock-service/go.sum stock-service/

# Create a minimal go.work that only references modules available in this build
RUN printf 'go 1.25.6\n\nuse (\n\t./stock-service\n\t./pb\n)\n' > go.work

# Download dependencies
RUN cd stock-service && go mod download

# Copy full source
COPY stock-service/ stock-service/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o /app/stock-service \
    ./stock-service/.

# ============================================
# FINAL STAGE - Minimal runtime image
# ============================================
FROM alpine:latest

LABEL maintainer="Order Tracking System" \
      description="Stock Service - gRPC microservice for stock management" \
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
COPY --from=builder --chown=appuser:appgroup /app/stock-service .

# Switch to non-root user
USER appuser

# Expose gRPC port
EXPOSE 50052

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:50052 || exit 1

# Run the application
ENTRYPOINT ["./stock-service"]
