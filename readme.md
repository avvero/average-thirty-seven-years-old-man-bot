# Average Thirty-Seven Years Old Man (bot)

Piece of digital art, telegram bot which acts like average thirty-seven years old man (by the author opinion):
- doesn't care about it at all
- tries to be humorous
- dumb
- arrogant
- spontaneous
- answers out of place
- uses bad word
- can understand toxicity

## Skills

Different ways of answering does call here as "skills".

| Name             | Probability (chance) | Explanation                                                                                           |
| ---------------- | -------------------- | ----------------------------------------------------------------------------------------------------- |
| Senseless phrase | 1 out of 100         | Random senseless phrase from memory                                                                   |
| Khaleesification | 1 out of 100         | Responses with mocking                                                                                |
| Huification      | 1 out of 50          | Responses with huifaed original message                                                               |
| gg               | 100 out of 100       | Responses: "gg" on the message with "gg" payload                                                      |
| morrowind        | 100 out of 100       | Responses: "Morrowind - одна из лучших игр эва" on the message with mention of this great game        |
| elden ring       | 100 out of 100       | Responses: "Elden Ring - это величие" on the message with mention of this great game                  |

... and many others.

## Run tests

```bash
go test ./...
```
## View test coverage

```bash
go test ./... -coverprofile cp.out
go tool cover -func cp.out | grep total | awk '{print $3}'
```

## Build

```bash
go build -o ./main ./cmd/bot/main.go
```

## Launch

```bash
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
require github.com/go -telegram-bot-api/telegram-bot-api/v5 v5.5.1
```

## TODO 

1. Improve build for docker