FROM golang:1.21-alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build both the main application and migration tool
RUN go build -o telemetry-ingest ./cmd/server && \
    go build -o migrate ./cmd/migrate

# Create a non-root user
RUN adduser -D -g '' appuser

# Switch to non-root user
USER appuser

# The default command will be overridden by docker-compose
CMD ["./telemetry-ingest"]

# Expose UDP port
EXPOSE 8089/udp 