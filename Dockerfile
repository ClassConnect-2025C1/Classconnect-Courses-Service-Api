FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application with specific flags for compatibility
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Use a smaller base image
FROM alpine:latest

WORKDIR /app

# Install CA certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main ./

# Copy .env file - fixed syntax
COPY .env ./

# Expose the port
EXPOSE 8002

# Command to run the application
CMD ["./main"]