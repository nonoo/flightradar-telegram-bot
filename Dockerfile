FROM golang:1.20 as builder
WORKDIR /app/
COPY go.mod go.sum /app/
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v

FROM alpine
COPY --from=builder /app/flightradar-telegram-bot /app/flightradar-telegram-bot

ENTRYPOINT ["/app/flightradar-telegram-bot"]
ENV BOT_TOKEN= ALLOWED_USERIDS= ADMIN_USERIDS= ALLOWED_GROUPIDS=
