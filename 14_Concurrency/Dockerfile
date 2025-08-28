# Step 1: Build stage
FROM golang:1.22 AS builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary (matches your Makefile)
RUN go build -o ./bin/server ./


# Step 2: Run stage
FROM alpine:latest

# Install certificates (needed if your Go app makes HTTPS requests)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/bin/server .

# Expose port (adjust to your app’s port)
EXPOSE 8080

# Run the binary
CMD ["./server"]
