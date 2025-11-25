# --- Stage 1: Build ---
# Use the official Go image as a builder image
FROM golang:1.19-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application for a Linux environment, statically linked
# CGO_ENABLED=0 is important for creating a truly static binary
# -ldflags="-w -s" strips debugging information, reducing the binary size
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o /main .

# --- Stage 2: Final ---
# Use a minimal image; scratch is the most minimal, but alpine is also a good choice if you need a shell
FROM scratch

# Copy the static binary from the builder stage
COPY --from=builder /main /main

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/main"]
