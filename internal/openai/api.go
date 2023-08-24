package openai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenAiClient struct {
	url       string
	accessKey string
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Choice []Choice `json:"choices"`
	Error  Error    `json:"error"`
}

type Choice struct {
	Message Message `json:"message"`
}

type Error struct {
	Message string `json:"message"`
}

func NewApiClient(url string, accessKey string) OpenAiClient {
	return OpenAiClient{url: url, accessKey: accessKey}
}

func (apiClient OpenAiClient) Completion(text string) (error, string) {
	return apiClient.CompletionByModel("gpt-3.5-turbo", text)
}

func (apiClient OpenAiClient) CompletionByModel(model string, text string) (error, string) {
	text = strings.ReplaceAll(text, "\"", "\\\"")

	fmt.Printf("Request to: %s\n with text: %s\n", apiClient.url, text)
	client := http.Client{Timeout: 500 * time.Second}

	requestBody, marshalError := json.Marshal(Request{
		Model: model,
		Messages: []Message{
			{Role: "user", Content: text},
		},
	})
	if marshalError != nil {
		return marshalError, ""
	}
	request, _ := http.NewRequest("POST", apiClient.url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiClient.accessKey)
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Response from: %s: %s", apiClient.url, err)
		return err, ""
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	fmt.Println("Response from openai: " + string(body))
	responseBody := &Response{}
	unmarshalError := json.Unmarshal(body, responseBody)
	if unmarshalError != nil {
		return errors.New(fmt.Sprintf("Response from: %s: %d: %s", apiClient.url, response.StatusCode, unmarshalError)), ""
	}
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Response from: %s: %d: %s", apiClient.url, response.StatusCode,
			responseBody.Error.Message)), ""
	}
	return err, responseBody.Choice[0].Message.Content
}
