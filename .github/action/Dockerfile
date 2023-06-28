FROM golang:1.20

# Make a directory for the files
RUN mkdir -p /app
WORKDIR /app

# Copy the files
COPY . ./

# Initialize go mod
RUN go mod init

# Build main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bruh ./cmd/bruh/main.go

# Make entrypoint.sh executable
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
