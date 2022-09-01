package huggingface

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

type HuggingFaceApiClient struct {
	acessKey string
}

type Labels struct {
	Title string  `json:"label"`
	Score float64 `json:"score"`
}

func NewHuggingFaceApiClient(acsessKey string) *HuggingFaceApiClient {
	return &HuggingFaceApiClient{acessKey: acsessKey}
}

func (apiClient *HuggingFaceApiClient) ToxicityScore(text string) (float64, error) {
	url := "https://api-inference.huggingface.co/models/apanc/russian-inappropriate-messages"
	fmt.Printf("Request to: %s\n", url)
	client := http.Client{Timeout: 5 * time.Second}

	requestBody, marshalError := json.Marshal(map[string]string{"inputs": text})
	if marshalError != nil {
		return 0, marshalError
	}
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiClient.acessKey)
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return 0, err
	}
	if response.StatusCode != 200 {
		fmt.Printf("Response code: %d\n", response.StatusCode)
		return 0, err
	}

	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	toxicityScore := 0.0
	toxicityFound := false

	bodyString := string(body)
	labels, err := Parse(bodyString)
	for _, label := range labels {
		if label.Title == "LABEL_1" {
			toxicityScore = label.Score
			toxicityFound = true
		}
	}
	if !toxicityFound {
		return 0, errors.New("Can't find LABEL_1 " + bodyString)
	}
	return toxicityScore, nil
}

func Parse(text string) ([]Labels, error) {
	text = strings.Replace(text, "[[", "[", -1)
	text = strings.Replace(text, "]]", "]", -1)

	var labels []Labels
	unmarshalError := json.Unmarshal([]byte(text), &labels)
	if unmarshalError != nil {
		return nil, unmarshalError
	} else {
		return labels, nil
	}
}
