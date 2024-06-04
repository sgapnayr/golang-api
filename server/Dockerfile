# Use the official Golang image to create a build artifact.
FROM golang:1.20-alpine AS builder

# Create and change to the /app directory.
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a minimal Docker image to reduce the size of the final image
FROM alpine:3.16

# Copy the binary from the builder stage
COPY --from=builder /app/main /main

# Command to run the executable
CMD ["/main"]
