package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/avvero/the_gamers_guild_bot/internal/data"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"github.com/avvero/the_gamers_guild_bot/internal/utils"
	"github.com/avvero/the_gamers_guild_bot/pkg/brain"
)

var (
	httpPort         = flag.String("httpPort", "8080", "http server port")
	token            = flag.String("token", "PROVIDE", "bot token")
	jsonBinMasterKey = flag.String("jsonBinMasterKey", "PROVIDE", "jsonBinMasterKey")
)

func main() {
	gracefullShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefullShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	flag.Parse()

	tokenEnv, found := os.LookupEnv("token")
	if found {
		token = &tokenEnv
	}
	jsonBinMasterKeyEnv, found := os.LookupEnv("json-bin-master-key")
	if found {
		jsonBinMasterKey = &jsonBinMasterKeyEnv
	}
	//Data
	jsonBinClient := data.NewJsonBinApiClient(*jsonBinMasterKey)
	data, err := jsonBinClient.Read()
	if err != nil {
		fmt.Printf("Could not read data: %s\n", err)
		panic(err)
	}
	ticker := time.NewTicker(1 * time.Minute)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Write data to bin", t)
				jsonBinClient.Write(data)
			}
		}
	}()
	scriber := statistics.NewScriber()
	brain := brain.NewBrain(brain.NewMemory(), true, scriber)

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"name\": \"Average Thirty-Seven Years Old Man (bot)\", \"version\": \"1.4\"}")
	})

	http.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("could not read body: %s\n", err)
			return
		}
		webhookRequest := &telegram.WebhookRequest{}
		json.Unmarshal(body, webhookRequest)
		if webhookRequest == nil || webhookRequest.Message == nil {
			fmt.Printf("could not unmarshal body\n")
			return
		}
		fmt.Println("Message from " + strconv.FormatInt(webhookRequest.Message.Chat.Id, 10) + " " +
			webhookRequest.Message.Chat.Title + ": " + webhookRequest.Message.Text)
		scriber.Keep(webhookRequest.Message)
		respond, response := brain.Decision(webhookRequest.Message.Chat.Id, webhookRequest.Message.Text)
		if respond {
			go func() {
				time.Sleep(time.Duration(utils.RandomUpTo(15)) * time.Second)
				sendMessage(webhookRequest.Message.Chat.Id, webhookRequest.Message.MessageId, response)
			}()
		}
	})

	log.Println("Http server started on port " + *httpPort)
	sendMessage(245851441, 0, "Bot is started, version 1.4")
	http.ListenAndServe(":"+*httpPort, nil)
	<-gracefullShutdown
	sendMessage(245851441, 0, "Bot is stopped, version 1.4")
	jsonBinClient.Write(data)
}

func sendMessage(chatId int64, receivedMessageId int64, message string) {
	replyToMessageId := ""
	if receivedMessageId != 0 {
		replyToMessageId = strconv.FormatInt(receivedMessageId, 10)
	}
	requestBody, marshalError := json.Marshal(map[string]string{
		"reply_to_message_id": replyToMessageId,
		"chat_id":             strconv.FormatInt(chatId, 10),
		"text":                message,
	})
	if marshalError != nil {
		fmt.Printf("could not marshal body: %s\n", marshalError)
	}
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
