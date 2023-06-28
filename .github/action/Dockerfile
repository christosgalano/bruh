# Use the official Go image as the base
FROM golang:1.20

# Create a directory for the files
RUN mkdir -p /app

# Copy the project files to the /app directory
COPY . /app

# Run go mod download from the /app directory
RUN cd /app && go mod download

# Build main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bruh /app/cmd/bruh/main.go

# Make entrypoint.sh executable
RUN chmod +x /app/entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]
