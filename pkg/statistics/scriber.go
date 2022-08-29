package statistics

import (
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"sync"
)

type Scriber struct {
	messages   chan *telegram.WebhookRequestMessage
	mutex      sync.Mutex
	statistics *Statistics
}

type Statistics struct {
	UsersStatistics map[string]*UserStatistics `json:"userStatistics"`
}

type UserStatistics struct {
	Username       string `json:"username"`
	MessageCounter int    `json:"messageCounter"`
}

func NewScriber() *Scriber {
	holder := &Scriber{
		messages:   make(chan *telegram.WebhookRequestMessage, 100),
		statistics: &Statistics{UsersStatistics: make(map[string]*UserStatistics)},
	}
	go holder.process()
	return holder
}

func (scriber Scriber) Keep(message *telegram.WebhookRequestMessage) {
	scriber.messages <- message
}

func (scriber Scriber) process() {
	for {
		select {
		case message := <-scriber.messages:
			scriber.mutex.Lock()
			userStatistics := scriber.statistics.UsersStatistics[message.From.Username]
			if userStatistics == nil {
				userStatistics = &UserStatistics{Username: message.From.Username}
				scriber.statistics.UsersStatistics[message.From.Username] = userStatistics
			}
			userStatistics.MessageCounter++
			scriber.mutex.Unlock()
		}
	}
}

func (scriber Scriber) GetStatistics() *Statistics {
	return scriber.statistics
}
