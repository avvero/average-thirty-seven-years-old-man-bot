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
	url       string
	accessKey string
}

type Labels struct {
	Title string  `json:"label"`
	Score float64 `json:"score"`
}

type ErrorDetails struct {
	Error         string  `json:"error"`
	EstimatedTime float64 `json:"estimated_time"`
}

func NewHuggingFaceApiClient(url string, accessKey string) *HuggingFaceApiClient {
	return &HuggingFaceApiClient{url: url, accessKey: accessKey}
}

func (apiClient *HuggingFaceApiClient) ToxicityScore(text string) (float64, error) {
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
		fmt.Printf("Request error: %s\n", err)
		return 0, err
	} else {
		fmt.Printf("Response code %d from : %s\n", response.StatusCode, apiClient.url)
	}
	if response.StatusCode == 503 {
		defer response.Body.Close()
		body, _ := io.ReadAll(response.Body)
		errorDetails := &ErrorDetails{}
		unmarshalError := json.Unmarshal(body, errorDetails)
		if unmarshalError != nil {
			return 0, unmarshalError
		}
		return 0, errors.New(fmt.Sprintf("Response code 503: %s, estimated time: %f", errorDetails.Error, errorDetails.EstimatedTime))
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
	fmt.Printf("Response from : %s :%s\n", apiClient.url, bodyString)
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
