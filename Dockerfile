FROM golang:latest
MAINTAINER avvero

ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/bot/main.go

ENTRYPOINT [ "/app/main" ]