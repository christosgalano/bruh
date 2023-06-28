FROM golang:1.20

# Make a directory for the files
RUN mkdir -p /app

# Copy the files
COPY . /app/

# Build main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bruh /app/cmd/bruh/main.go

# Make entrypoint.sh executable
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
