# The GamersGuild Bot

## Launch
```bash
go build
./main -httpPort=8080
```

## CI
```bash
docker-compose down --rmi all && docker-compose up -d
```

## https://go-telegram-bot-api.dev/getting-started/index.html

Let's use: https://go-telegram-bot-api.dev/getting-started/index.html

In `go.mod`
```go
require github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
```