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
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	Temperature      float64 `json:"temperature"`
	MaxTokens        int     `json:"max_tokens"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  float64 `json:"presence_penalty"`
}

type Response struct {
	Choice []Choice `json:"choices"`
	Error  Error    `json:"error"`
}

type Choice struct {
	Text string `json:"text"`
}

type Error struct {
	Message string `json:"message"`
}

func NewApiClient(url string, accessKey string) OpenAiClient {
	return OpenAiClient{url: url, accessKey: accessKey}
}

func (apiClient OpenAiClient) Completion(text string) (error, string) {
	fmt.Printf("Request to: %s\n", apiClient.url)
	client := http.Client{Timeout: 500 * time.Second}

	text = strings.ReplaceAll(text, "\"", "\\\"")

	requestBody, marshalError := json.Marshal(Request{Model: "text-davinci-003", Prompt: text, Temperature: 0.9,
		MaxTokens: 4097 - len(text), TopP: 1, FrequencyPenalty: 0.0, PresencePenalty: 0.6})
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
	responseBody := &Response{}
	unmarshalError := json.Unmarshal(body, responseBody)
	if unmarshalError != nil {
		return errors.New(fmt.Sprintf("Response from: %s: %d: %s", apiClient.url, response.StatusCode, unmarshalError)), ""
	}
	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Response from: %s: %d: %s", apiClient.url, response.StatusCode,
			responseBody.Error.Message)), ""
	}
	result := strings.ReplaceAll(responseBody.Choice[0].Text, "\n\n", "")
	result = strings.ReplaceAll(result, "\"", "")
	return err, result
}
