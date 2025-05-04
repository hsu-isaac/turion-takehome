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

# Build the application
RUN go build -o telemetry-ingest ./cmd/server

# Run the application
CMD ["./telemetry-ingest"]

# Expose UDP port
EXPOSE 8089/udp 