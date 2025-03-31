# Use a minimal Go base image
FROM golang:1.22.5 AS builder

# Set working directory
WORKDIR /app

# Copy Go files
COPY . .

# Download dependencies
RUN go mod tidy

# Build the Go binary
RUN go build -o app .

# Use a lightweight image for running the app
FROM alpine:latest

# Install SQLite and CA certificates
RUN apk add --no-cache sqlite

# Set working directory
WORKDIR /app

# Copy built binary from the builder stage
COPY --from=builder /app/app .

# Run the application
CMD ["./app"]

