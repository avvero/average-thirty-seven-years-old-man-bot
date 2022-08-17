# Average Thirty-Seven Years Old Man (bot)

Piece of digital art, telegram bot which acts like average thirty-seven years old man (by the author opinion):

- doesn't give a shit
- tries to be humorous
- dumb
- arrogant
- spontaneous
- answers out of place
- uses bad words

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
require github.com/go -telegram-bot-api/telegram-bot-api/v5 v5.5.1
```
