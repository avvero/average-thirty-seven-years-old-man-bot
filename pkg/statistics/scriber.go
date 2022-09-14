package statistics

import (
	"github.com/avvero/the_gamers_guild_bot/internal/data"
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
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
			// By user
			chatStatistics := scriber.data.ChatStatistics[message.Chat.Id]
			if chatStatistics == nil {
				chatStatistics = &data.ChatStatistics{UsersStatistics: make(map[string]*data.MessageStatistics)}
				scriber.data.ChatStatistics[message.Chat.Id] = chatStatistics
			}
			user := message.From.Username
			if user == "" {
				user = message.From.LastName + " " + message.From.FirstName
			}
			userStatistics := chatStatistics.UsersStatistics[user]
			if userStatistics == nil {
				userStatistics = &data.MessageStatistics{}
				chatStatistics.UsersStatistics[user] = userStatistics
			}
			userStatistics.MessageCounter++
			// Daily
			date := time.Now().Format("2006-01-02")
			if chatStatistics.DailyStatistics == nil {
				chatStatistics.DailyStatistics = make(map[string]*data.MessageStatistics)
			}
			dailyStatistics := chatStatistics.DailyStatistics[date]
			if dailyStatistics == nil {
				dailyStatistics = &data.MessageStatistics{}
				chatStatistics.DailyStatistics[date] = dailyStatistics
			}
			dailyStatistics.MessageCounter++
			scriber.mutex.Unlock()
		}
	}
}

func (scriber Scriber) GetStatistics(chatId int64) *data.ChatStatistics {
	return scriber.data.ChatStatistics[chatId]
}

// TODO move to external object
func (scriber Scriber) GetStatisticsPrettyPrint(chatId int64) string {
	if scriber.data.ChatStatistics == nil ||
		scriber.data.ChatStatistics[chatId] == nil ||
		scriber.data.ChatStatistics[chatId].UsersStatistics == nil ||
		scriber.data.ChatStatistics[chatId].DailyStatistics == nil {
		return ""
	}
	chatStatistics := scriber.data.ChatStatistics[chatId]

	sb := strings.Builder{}
	sb.WriteString("Statistics by user:\n")
	usKeys := sortedKeys(chatStatistics.UsersStatistics)
	for _, k := range usKeys {
		sb.WriteString(" - " + k + ": " + strconv.Itoa(chatStatistics.UsersStatistics[k].MessageCounter) + "\n")
	}
	sb.WriteString("Statistics by day:\n")
	dsKeys := sortedKeys(chatStatistics.DailyStatistics)
	for _, k := range dsKeys {
		sb.WriteString(" - " + k + ": " + strconv.Itoa(chatStatistics.DailyStatistics[k].MessageCounter) + "\n")
	}
	return sb.String()
}

func sortedKeys(m map[string]*data.MessageStatistics) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
