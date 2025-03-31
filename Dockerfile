# Use a minimal Go base image
FROM golang:1.21 AS builder

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

# Copy SQLite database file if you have a pre-existing one
COPY recipes.db .

# Run the application
CMD ["./app"]

