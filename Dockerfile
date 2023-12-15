# Use a builder image to build the Go project
FROM golang:alpine AS builder

# Set the working directory inside the builder container
WORKDIR /app
# Copy the Go project files to the builder container
COPY . .
# Build the Go project
RUN go build -o auth

# Use a minimal base image for the final image
FROM alpine
# Set the working directory inside the final container
WORKDIR /app
# Copy the compiled binary from the builder container to the final container
COPY --from=builder /app/auth .
# Run the compiled binary
CMD ["./auth"]
