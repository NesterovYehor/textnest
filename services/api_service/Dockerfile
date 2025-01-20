# Use a lightweight Go image as the base for the build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies (e.g., git) in the builder stage
RUN apk add --no-cache git

ENV GOFLAGS=-mod=readonly

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download and verify dependencies
RUN go mod download && go mod verify

# Copy the entire project source code
COPY . .

# Build the Go application and create a binary named `app`
RUN go build -o /app/build/app ./cmd/main.go

# Use a minimal image for the runtime stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/build/app .

# Run the binary
CMD ["./app"]

