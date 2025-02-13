# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk update && apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./main.go

# Stage 2: Run the binary in a minimal container
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
