FROM golang:alpine AS builder

ARG CGO_ENABLED=0
ARG GOOS=linux

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal

RUN go build -o app ./cmd/api

FROM alpine
LABEL org.opencontainers.image.source=https://github.com/Eretic431/datingTelegramBot

WORKDIR /app

COPY --from=builder /build/app .
EXPOSE 80
CMD ["./app"]