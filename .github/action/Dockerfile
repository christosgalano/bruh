FROM golang:1.20

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/bruh ./cmd/bruh/main.go

RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]
