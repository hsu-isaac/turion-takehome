FROM golang:1.21-alpine

WORKDIR /app

# Copy go mod files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o telemetry-generator ./cmd/server

# Run the application
CMD ["./telemetry-generator"] 