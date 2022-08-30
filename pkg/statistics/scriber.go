package statistics

import (
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"strconv"
	"sync"
)

type Scriber struct {
	messages   chan *telegram.WebhookRequestMessage
	mutex      sync.Mutex
	statistics *Statistics
}

type Statistics struct {
	ChatStatistics map[string]*ChatStatistics `json:"chatStatistics"`
}

type ChatStatistics struct {
	UsersStatistics map[string]*UserStatistics `json:"userStatistics"`
}

type UserStatistics struct {
	Username       string `json:"username"`
	MessageCounter int    `json:"messageCounter"`
}

func NewScriber() *Scriber {
	holder := &Scriber{
		messages:   make(chan *telegram.WebhookRequestMessage, 100),
		statistics: &Statistics{ChatStatistics: make(map[string]*ChatStatistics)},
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
			chatId := strconv.FormatInt(message.Chat.Id, 10)
			chatStatistics := scriber.statistics.ChatStatistics[chatId]
			if chatStatistics == nil {
				chatStatistics = &ChatStatistics{UsersStatistics: make(map[string]*UserStatistics)}
				scriber.statistics.ChatStatistics[chatId] = chatStatistics
			}
			userStatistics := chatStatistics.UsersStatistics[message.From.Username]
			if userStatistics == nil {
				userStatistics = &UserStatistics{Username: message.From.Username}
				chatStatistics.UsersStatistics[message.From.Username] = userStatistics
			}
			// Set
			userStatistics.MessageCounter++
			scriber.mutex.Unlock()
		}
	}
}

func (scriber Scriber) GetStatistics(chatId int64) *ChatStatistics {
	return scriber.statistics.ChatStatistics[strconv.FormatInt(chatId, 10)]
}
