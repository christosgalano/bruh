FROM golang:1.20

# Make a directory for the files
RUN mkdir -p /app

# Copy the files
COPY internal /app/internal
COPY cmd/ /app/cmd
COPY go.mod go.sum /app/
COPY entrypoint.sh /app/

# Install dependencies
RUN cd /app/ && go mod download

# Build main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bruh /app/cmd/bruh/main.go

# Make entrypoint.sh executable
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
