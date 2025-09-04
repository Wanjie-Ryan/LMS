# Stage 1: Build the Go application
# use official Go image with version 1.23.2 to compile code
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
# the directory exists inside the image.
WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
# download all go dependencies into the image
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary (main.go) and name the file lms-server
# the -o allows one to name the binary output, that is genearted after being built
# RUN go build -o lms-server main.go
# RUN go build -o lms-server ./cmd/api/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lms-server ./cmd/api/


# Stage 2: Create minimal image to run the binary
FROM debian:bullseye-slim

# Set working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/lms-server .
# Copy the .env file into the container
COPY --from=builder /app/.env .env

# Expose the port your app runs on
EXPOSE 8080

# Start the server
CMD ["./lms-server"]
