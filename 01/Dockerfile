# Build Stage
FROM golang:1.21.4-alpine AS build

WORKDIR /src

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Enable CGO so SQLite works
ENV CGO_ENABLED=1

# Build the Go binary
RUN go build -o server ./cmd/api

# Final Stage (MUST match glibc version)
FROM debian:bookworm-slim  
# ✅ This is where your binary will run

# Install runtime dependencies (e.g., for SQLite)
RUN apt-get update && apt-get install -y ca-certificates libsqlite3-0 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# ✅ Create a 'storage' directory
RUN mkdir -p /app/storage


# Copy the binary and config from the build stage
COPY --from=build /src/server .
COPY config/ ./config/

EXPOSE 7000

# Run the server
ENTRYPOINT ["./server", "-config", "./config/local.yaml"]
