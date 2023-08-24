package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type WebhookRequest struct {
	Message *WebhookRequestMessage `json:"message"`
}

type WebhookRequestMessage struct {
	MessageId            int64                        `json:"message_id"`
	From                 *WebhookRequestMessageSender `json:"from"`
	Chat                 *WebhookRequestMessageChat   `json:"chat"`
	Text                 string                       `json:"text"`
	ForwardFromMessageId int64                        `json:"forward_from_message_id"`
	NewChatParticipant   *NewChatParticipant          `json:"new_chat_participant"`
}

type WebhookRequestMessageSender struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type NewChatParticipant struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

type WebhookRequestMessageChat struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

type TelegramApiClient struct {
	host  string
	token string
}

func NewTelegramApiClient(host string, token string) *TelegramApiClient {
	return &TelegramApiClient{host: host, token: token}
}

func (apiClient TelegramApiClient) SendMessage(chatId int64, receivedMessageId int64, message string) {
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
	url := apiClient.host + "/bot" + apiClient.token + "/sendMessage"
	fmt.Printf("Request to: %s, message: %s\n", url, requestBody)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("could not read response body: %s\n", err)
		return
	}
	fmt.Println("Telegram Response body " + string(body))
	//	webhookRequest := &telegram.WebhookRequest{}
	//	json.Unmarshal(body, webhookRequest)
}

func (apiClient TelegramApiClient) SendMessage2(chatId int64, receivedMessageId int64, message string) {
	replyToMessageId := ""
	if receivedMessageId != 0 {
		replyToMessageId = strconv.FormatInt(receivedMessageId, 10)
	}
	requestBody, marshalError := json.Marshal(map[string]string{
		"reply_to_message_id": replyToMessageId,
		"chat_id":             strconv.FormatInt(chatId, 10),
		"text":                message,
		"parse_mode":          "markdown",
	})
	if marshalError != nil {
		fmt.Printf("could not marshal body: %s\n", marshalError)
	}
	client := http.Client{Timeout: 5 * time.Second}
	url := apiClient.host + "/bot" + apiClient.token + "/sendMessage"
	fmt.Printf("Request to: %s, message: %s\n", url, message)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	_, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return
	}
}

func (apiClient TelegramApiClient) SendSticker(chatId int64, receivedMessageId int64, fileId string) {
	replyToMessageId := ""
	if receivedMessageId != 0 {
		replyToMessageId = strconv.FormatInt(receivedMessageId, 10)
	}
	requestBody, marshalError := json.Marshal(map[string]string{
		"reply_to_message_id": replyToMessageId,
		"chat_id":             strconv.FormatInt(chatId, 10),
		"sticker":             fileId,
	})
	if marshalError != nil {
		fmt.Printf("could not marshal body: %s\n", marshalError)
	}
	client := http.Client{Timeout: 5 * time.Second}
	url := apiClient.host + "/bot" + apiClient.token + "/sendSticker"
	fmt.Printf("Request to: %s, sticker: %s\n", url, fileId)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	_, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return
	}
}

func (apiClient TelegramApiClient) BanChatMember(chatId int64, userId int64) {
	requestBody, marshalError := json.Marshal(map[string]string{
		"chat_id": strconv.FormatInt(chatId, 10),
		"user_id": strconv.FormatInt(userId, 10),
	})
	if marshalError != nil {
		fmt.Printf("could not marshal body: %s\n", marshalError)
	}
	client := http.Client{Timeout: 5 * time.Second}
	url := apiClient.host + "/bot" + apiClient.token + "/banChatMember"
	fmt.Printf("Request to: %s, user: %s\n", url, strconv.FormatInt(userId, 10))
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	_, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return
	}
}

func (apiClient TelegramApiClient) SendChatAction(chatId int64, action string) {
	requestBody, marshalError := json.Marshal(map[string]string{
		"chat_id": strconv.FormatInt(chatId, 10),
		"action":  action,
	})
	if marshalError != nil {
		fmt.Printf("could not marshal body: %s\n", marshalError)
	}
	client := http.Client{Timeout: 5 * time.Second}
	url := apiClient.host + "/bot" + apiClient.token + "/sendChatAction"
	fmt.Printf("Request to: %s, action: %s\n", url, action)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	_, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return
	}
}
