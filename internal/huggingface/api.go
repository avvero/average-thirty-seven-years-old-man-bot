package huggingface

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HuggingFaceApiClient struct {
	url       string
	accessKey string
}

type Label struct {
	Title string  `json:"label"`
	Score float64 `json:"score"`
}

type ErrorDetails struct {
	Error         string  `json:"error"`
	EstimatedTime float64 `json:"estimated_time"`
}

func NewApiClient(url string, accessKey string) HuggingFaceApiClient {
	return HuggingFaceApiClient{url: url, accessKey: accessKey}
}

func (apiClient HuggingFaceApiClient) ToxicityScore(text string) (float64, error) {
	fmt.Printf("Request to: %s\n", apiClient.url)
	client := http.Client{Timeout: 5 * time.Second}

	requestBody, marshalError := json.Marshal(map[string]string{"inputs": text})
	if marshalError != nil {
		return 0, marshalError
	}
	request, _ := http.NewRequest("POST", apiClient.url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+apiClient.accessKey)
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Response from: %s: %s", apiClient.url, err)
		return 0, err
	}
	if response.StatusCode == 503 {
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		errorDetails := &ErrorDetails{}
		unmarshalError := json.Unmarshal(body, errorDetails)
		if unmarshalError != nil {
			return 0, unmarshalError
		}
		return 0, errors.New(fmt.Sprintf("Response from: %s: 503: %s, estimated time: %f", apiClient.url,
			errorDetails.Error, errorDetails.EstimatedTime))
	}
	if response.StatusCode != 200 {
		return 0, errors.New(fmt.Sprintf("Response from: %s: %d", apiClient.url, response.StatusCode))
	}

	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	var labels [][]Label
	unmarshalError := json.Unmarshal(body, &labels)
	if unmarshalError != nil {
		return 0, errors.New(fmt.Sprintf("Can't parse response: %s: %s", string(body), unmarshalError))
	}
	fmt.Printf("Response from: %s: %d: %v\n", apiClient.url, response.StatusCode, labels)
	// Arrays composition
	for _, labels2 := range labels {
		for _, label := range labels2 {
			if label.Title == "LABEL_1" {
				return label.Score, nil
			}
		}
	}
	return 0, errors.New(fmt.Sprintf("Can't find LABEL_1: %v", labels))
}
