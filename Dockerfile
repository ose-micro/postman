# Use official Go base image
FROM golang:1.24.3-alpine AS builder

# Enable Go modules and CGO disabled for static binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

# Install git and clean up cache
RUN apk add --no-cache git && apk --no-cache upgrade

# Copy go files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source
COPY . .

# Build the binary
RUN go build -o ose-postman ./cmd

# ─────────────────────────────

# Final lightweight image
FROM alpine:latest

# Add necessary packages for ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/ose-postman .

# Set timezone and environment variables if needed
ENV TZ=UTC

# Command
ENTRYPOINT ["./postman"]
