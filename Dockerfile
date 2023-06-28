FROM golang:1.20

# Make a directory for the files
RUN mkdir -p /app

# Install dependencies
COPY go.mod go.sum /app/
RUN cd /app/ && go mod download

# Copy and build main.go
COPY cmd/bruh/main.go /app/
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bruh /app/main.go

# Copy and make entrypoint.sh executable
COPY entrypoint.sh /app/
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
