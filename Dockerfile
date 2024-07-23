# Multi-stage build for Golang application
FROM golang:1.22.3 AS builder

# Set working directory
WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the Golang application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o go-skeleton .

# Build the final image using Alpine
FROM alpine:latest

# Set the timezone to Asia/Jakarta
ENV TZ=Asia/Jakarta
RUN apk add -U tzdata

# Set working directory for the final image
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/go-skeleton .

# Switch to the application user to run the application
# This is a recommended security practice to run the application as a non-root user
# By default, the application is run as the user "app" defined in the Dockerfile
# This is a good practice to avoid potential security risks
USER    app

# Expose port 8000
EXPOSE 8000

# Command to run the application
CMD ["./go-skeleton"]
