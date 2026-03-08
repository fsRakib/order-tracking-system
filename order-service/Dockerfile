# ============================================
# BUILD STAGE - Compile Go application
# ============================================
FROM golang:1.25.6-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy the shared pb module
COPY pb/go.mod pb/go.sum pb/
COPY pb/order/ pb/order/
COPY pb/stock/ pb/stock/

# Copy the order-service module
COPY order-service/go.mod order-service/go.sum order-service/

# Create a minimal go.work that only references modules available in this build
RUN printf 'go 1.25.6\n\nuse (\n\t./order-service\n\t./pb\n)\n' > go.work

# Download dependencies
RUN cd order-service && go mod download

# Copy full source
COPY order-service/ order-service/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o /app/order-service \
    ./order-service/.

# ============================================
# FINAL STAGE - Minimal runtime image
# ============================================
FROM alpine:latest

LABEL maintainer="Order Tracking System" \
      description="Order Service - gRPC microservice for order management" \
      version="1.0"

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata netcat-openbsd && \
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

# Health check - TCP probe (gRPC does not speak plain HTTP)
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD nc -z localhost 50051 || exit 1

# Run the application
ENTRYPOINT ["./order-service"]
