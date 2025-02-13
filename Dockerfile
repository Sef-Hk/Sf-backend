# Stage 1: Build the Go binary
FROM golang:1.19-alpine AS builder

# Set the working directory
WORKDIR /app

# Install CA certificates and Git (needed for HTTPS and module fetching)
RUN apk update && apk add --no-cache ca-certificates git

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o app ./main.go

# Stage 2: Run the binary in a minimal container
FROM alpine:latest
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Expose the port your app listens on (adjust if needed)
EXPOSE 8080

# Run the binary
CMD ["./app"]
