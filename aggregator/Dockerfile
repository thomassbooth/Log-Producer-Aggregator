# Use the official Golang image as a build stage
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files from the root directory
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container from the aggregator directory
COPY aggregator/ ./aggregator/

# Build the Go app from the cmd directory
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./aggregator/cmd

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8005

# Command to run the executable
CMD ["./main"]
