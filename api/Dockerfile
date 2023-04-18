# Use the official Golang base image
FROM golang:1.20 as builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY ./ ./

# Download the Go module dependencies
RUN go mod download

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o app api/main.go

# Create a minimal Docker image for the executable
FROM scratch
COPY --from=builder /app/app .
ENTRYPOINT ["./app"]
