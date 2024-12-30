FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main ./cmd/api

# Use a smaller base image for the final stage
FROM alpine:latest

WORKDIR /app

# Copy only the binary from builder
COPY --from=builder /build/main .

# Add ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"] 