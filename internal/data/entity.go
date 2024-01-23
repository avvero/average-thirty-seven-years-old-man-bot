package data

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Data struct {
	ChatStatistics map[int64]*ChatStatistics `json:"chatStatistics"`
}

type ChatStatistics struct {
	UsersStatistics     map[string]*UserMessageStatistics `json:"userStatistics"`
	HouseStatistics     map[string]*HouseStatistics       `json:"houseStatistics,omitempty"`
	DailyStatistics     map[string]*DayMessageStatistics  `json:"dailyStatistics"`
	DailyWordStatistics map[string]map[string]int         `json:"-"`
	Notifications       map[string]*Notification          `json:"notifications,omitempty"`
	Messages            []*Message                        `json:"-"`
}

type UserMessageStatistics struct {
	MessageCounter      int     `json:"messageCounter"`
	ToxicityScore       float64 `json:"toxicityScore"`
	Tension             int     `json:"tension"`
	LastMessageDate     string  `json:"lastMessageDate,omitempty"`
	LastMessageDateTime string  `json:"lastMessageDateTime,omitempty"`
}

type HouseStatistics struct {
	Score int `json:"messageCounter"`
}

type DayMessageStatistics struct {
	MessageCounter int     `json:"messageCounter"`
	ToxicityScore  float64 `json:"toxicityScore"`
}

type Notification struct {
	User   string
	Action string
	Time   string
}

type Message struct {
	User string
	Text string
}

func NewLocalStorage() LocalStorage {
	return LocalStorage{}
}

type LocalStorage struct {
	mutex sync.Mutex
}

func (storage LocalStorage) Read() (*Data, error) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	fileData, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Printf("File read error: %s\n", err)
		return nil, err
	}

	data := &Data{}
	unmarshalError := json.Unmarshal(fileData, data)
	if unmarshalError != nil {
		fmt.Printf("Could not read data: %s\n", unmarshalError)
		return nil, unmarshalError
	}

	return data, nil
}

func (storage LocalStorage) Write(data *Data) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Marshalling error: %s\n", err)
		return err
	}

	writeErr := os.WriteFile("data.json", jsonData, 0644)
	if writeErr != nil {
		fmt.Printf("File write error: %s\n", writeErr)
		return writeErr
	}

	return nil
}
