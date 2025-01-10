# Use the official Go image as the base
FROM golang:1.23

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
RUN go build -o build/app ./cmd/main.go

CMD ["build/app"]


