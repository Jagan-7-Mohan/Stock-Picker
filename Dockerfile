# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code (excluding files in .dockerignore)
COPY . .

# Verify the cmd directory exists
RUN ls -la /app/cmd/stock-picker/ || (echo "ERROR: cmd/stock-picker directory not found!" && exit 1)

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o stock-picker ./cmd/stock-picker

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and tzdata for timezone support
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/stock-picker .

# Set timezone to IST
ENV TZ=Asia/Kolkata

# Run the application
CMD ["./stock-picker"]

