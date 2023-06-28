FROM golang:1.20

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal
COPY entrypoint.sh ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bruh ./cmd/bruh/main.go

RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
