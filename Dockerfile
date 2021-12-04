FROM golang:1.17 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api

FROM alpine:latest
LABEL org.opencontainers.image.source=https://github.com/Eretic431/datingTelegramBot

WORKDIR /app

COPY --from=builder /build/app .
EXPOSE 80
CMD ["./app"]