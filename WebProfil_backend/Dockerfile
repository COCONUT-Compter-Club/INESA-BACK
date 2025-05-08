# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o godesaapps-main .

# Stage 2: Run the application
FROM alpine:latest

WORKDIR /app

# Copy the binary and necessary files from builder
COPY --from=builder /app/godesaapps-main .
COPY --from=builder /app/.env .

# Expose the application port
EXPOSE 8085

# Command to run the application
CMD ["./godesaapps-main"]
