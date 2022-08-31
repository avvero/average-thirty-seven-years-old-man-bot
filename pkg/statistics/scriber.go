package statistics

import (
	"github.com/avvero/the_gamers_guild_bot/internal/data"
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"sync"
)

type Scriber struct {
	messages chan *telegram.WebhookRequestMessage
	mutex    sync.Mutex
	data     *data.Data
}

func NewScriberWithData(data *data.Data) *Scriber {
	holder := &Scriber{
		messages: make(chan *telegram.WebhookRequestMessage, 100),
		data:     data,
	}
	go holder.process()
	return holder
}

func NewScriber() *Scriber {
	return NewScriberWithData(&data.Data{ChatStatistics: make(map[int64]*data.ChatStatistics)})
}

func (scriber Scriber) Keep(message *telegram.WebhookRequestMessage) {
	scriber.messages <- message
}

func (scriber Scriber) process() {
	for {
		select {
		case message := <-scriber.messages:
			scriber.mutex.Lock()
			chatStatistics := scriber.data.ChatStatistics[message.Chat.Id]
			if chatStatistics == nil {
				chatStatistics = &data.ChatStatistics{UsersStatistics: make(map[string]*data.UserStatistics)}
				scriber.data.ChatStatistics[message.Chat.Id] = chatStatistics
			}
			user := message.From.Username
			if user == "" {
				user = message.From.LastName + " " + message.From.FirstName
			}
			userStatistics := chatStatistics.UsersStatistics[user]
			if userStatistics == nil {
				userStatistics = &data.UserStatistics{Username: user}
				chatStatistics.UsersStatistics[user] = userStatistics
			}
			// Set
			userStatistics.MessageCounter++
			scriber.mutex.Unlock()
		}
	}
}

func (scriber Scriber) GetStatistics(chatId int64) *data.ChatStatistics {
	return scriber.data.ChatStatistics[chatId]
}
