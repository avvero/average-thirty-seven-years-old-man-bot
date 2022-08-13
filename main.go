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
	"time"
)

var (
	httpPort = flag.String("httpPort", "8080", "http server port")
	token    = flag.String("token", "PROVIDE", "bot token")
)

func main() {
	flag.Parse()

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"name\": \"TheGamerGuildBot\"}")
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
		fmt.Println("Chat message: " + webhookRequest.Message.Text)
		if webhookRequest.Message.Text == "gg" {
			requestBody, _ := json.Marshal(map[string]string{
				"chat_id": "-1001733786877",
				"text":    "gg",
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
	})

	log.Println("Http server started on port " + *httpPort)
	http.ListenAndServe(":"+*httpPort, nil)
}

type WebhookRequest struct {
	Message *WebhookRequestMessage `json:"message"`
}

type WebhookRequestMessage struct {
	Text string `json:"text"`
}
