# ----- Build Stage -----
FROM golang:1.24.1-alpine AS builder

# Install necessary tools
RUN apk add --no-cache git

WORKDIR /app

# Copy Go modules files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the app
RUN go build -o ose-postman ./cmd

# ----- Run Stage -----
FROM alpine:3.18

# Install certs (for HTTPS, Mongo, etc.)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/ose-postman .
COPY --from=builder /app/cmd/config.production.json ./cmd/config.production.json

# Expose the gRPC port (adjust if needed)
EXPOSE 50051

# Run the service
CMD ["./ose-postman"]
