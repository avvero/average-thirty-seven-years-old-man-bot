package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Data struct {
	ChatStatistics map[int64]*ChatStatistics `json:"chatStatistics"`
}

type ChatStatistics struct {
	UsersStatistics     map[string]*MessageStatistics `json:"userStatistics"`
	DailyStatistics     map[string]*MessageStatistics `json:"dailyStatistics"`
	DailyWordStatistics map[string]map[string]int     `json:"-"`
}

type MessageStatistics struct {
	MessageCounter int     `json:"messageCounter"`
	ToxicityScore  float64 `json:"toxicityScore"`
}

type JsonBinResponse struct {
	Record *Data `json:"record"`
}

type JsonBinApiClient struct {
	binId     string
	masterKey string
}

func NewJsonBinApiClient(masterKey string) *JsonBinApiClient {
	return &JsonBinApiClient{binId: "6377909e2b3499323b039787", masterKey: masterKey}
}

func (apiClient *JsonBinApiClient) Read() (*Data, error) {
	client := http.Client{Timeout: 5 * time.Second}
	url := "https://api.jsonbin.io/v3/b/" + apiClient.binId
	fmt.Printf("Request to: %s\n", url)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Master-Key", apiClient.masterKey)
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return nil, err
	}
	if response.StatusCode != 200 {
		fmt.Printf("Response code: %d\n", response.StatusCode)
		return nil, err
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	jsonBinResponse := &JsonBinResponse{}
	unmarshalError := json.Unmarshal(body, jsonBinResponse)
	if unmarshalError != nil {
		fmt.Printf("Could not read data: %s\n", unmarshalError)
		return nil, unmarshalError
	}
	return jsonBinResponse.Record, nil
}

func (apiClient *JsonBinApiClient) Write(data *Data) error {
	requestBody, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Marshalling error: %s\n", err)
		return err
	}
	client := http.Client{Timeout: 5 * time.Second}
	url := "https://api.jsonbin.io/v3/b/" + apiClient.binId
	fmt.Printf("Request to: %s\n", url)
	request, _ := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Master-Key", apiClient.masterKey)
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return err
	}
	if response.StatusCode != 200 {
		fmt.Printf("Response code: %d\n", response.StatusCode)
		return err
	}
	return nil
}
