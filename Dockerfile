# ----- Build Stage -----
FROM golang:1.25 AS builder

# Ensure Go uses correct toolchain automatically
ENV GOTOOLCHAIN=auto

# Install necessary tools
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy Go modules files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build statically linked binary for Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o postman ./cmd

# ----- Run Stage -----
FROM alpine:3.18

ARG PORT=8080

# Install runtime dependencies
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/postman .
COPY --from=builder /app/cmd/config.development.json ./cmd/config.development.json
COPY --from=builder /app/cmd/config.production.json ./cmd/config.production.json

# Expose the gRPC port (adjust if needed)
ENV PORT=$PORT

# Run the service
CMD ["./postman"]