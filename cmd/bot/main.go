package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"

	"github.com/avvero/the_gamers_guild_bot/internal/utils"
	"github.com/avvero/the_gamers_guild_bot/pkg/brain"
)

var (
	httpPort = flag.String("httpPort", "8080", "http server port")
	token    = flag.String("token", "PROVIDE", "bot token")
)

func main() {
	flag.Parse()

	tokenEnv, found := os.LookupEnv("token")
	if found {
		token = &tokenEnv
	}
	brain := brain.NewBrain(brain.NewMemory(), true)

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"name\": \"TheGamerGuildBot\", \"version\": \"1.3\"}")
	})

	http.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("could not read body: %s\n", err)
			return
		}
		webhookRequest := &WebhookRequest{}
		json.Unmarshal(body, webhookRequest)
		if webhookRequest == nil || webhookRequest.Message == nil {
			fmt.Printf("could not unmarshal body\n")
			return
		}
		fmt.Println("Message from " + strconv.FormatInt(webhookRequest.Message.Chat.Id, 10) + " " +
			webhookRequest.Message.Chat.Title + ": " + webhookRequest.Message.Text)

		respond, response := brain.Decision(webhookRequest.Message.Chat.Id, webhookRequest.Message.Text)
		if respond {
			go func() {
				time.Sleep(time.Duration(utils.RandomUpTo(15)) * time.Second)
				sendMessage(webhookRequest.Message.Chat.Id, webhookRequest.Message.MessageId, response)
			}()
		}
	})

	log.Println("Http server started on port " + *httpPort)
	sendMessage(245851441, 0, "Bot is redeployed")
	http.ListenAndServe(":"+*httpPort, nil)
}

type WebhookRequest struct {
	Message *WebhookRequestMessage `json:"message"`
}

type WebhookRequestMessage struct {
	MessageId int64                        `json:"message_id"`
	From      *WebhookRequestMessageSender `json:"from"`
	Chat      *WebhookRequestMessageChat   `json:"chat"`
	Text      string                       `json:"text"`
}

type WebhookRequestMessageSender struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	username  string `json:"username"`
}

type WebhookRequestMessageChat struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

func sendMessage(chatId int64, receivedMessageId int64, message string) {
	replyToMessageId := ""
	if receivedMessageId != 0 {
		replyToMessageId = strconv.FormatInt(receivedMessageId, 10)
	}
	requestBody, _ := json.Marshal(map[string]string{
		"reply_to_message_id": replyToMessageId,
		"chat_id":             strconv.FormatInt(chatId, 10),
		"text":                message,
	})
	client := http.Client{Timeout: 5 * time.Second}
	url := "https://api.telegram.org/bot" + *token + "/sendMessage"
	fmt.Printf("Request to: %s\n", url)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	_, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return
	}
}
