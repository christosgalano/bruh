# Use the official Go image as the base
FROM golang:1.20

# Create a directory for the files
RUN mkdir -p /app

# Copy the project files to the /app directory
COPY . /app

# Change the current working directory to /app
WORKDIR /app

# Run go mod download
RUN go mod download

# Build main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o bruh ./cmd/bruh/main.go

# Make entrypoint.sh executable
RUN chmod +x entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]
