# Average Thirty-seven Years Old Man (bot)

Piece of digital art, telegram bot which acts like average thirty-seven years old man (by the author opinion):
- don't give a shit
- tries to be humorous 
- dumb
- arrogant
- spontaneous
- answers out of place
- uses bad words

## Run tests
```bash
go test ./...
```

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
