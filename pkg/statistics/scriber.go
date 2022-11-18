package statistics

import (
	"github.com/avvero/the_gamers_guild_bot/internal/data"
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Scriber struct {
	messages       chan *telegram.WebhookRequestMessage
	mutex          sync.Mutex
	data           *data.Data
	statisticsPage string
}

func NewScriberWithData(data *data.Data, statisticsPage string) *Scriber {
	holder := &Scriber{
		messages:       make(chan *telegram.WebhookRequestMessage, 100),
		data:           data,
		statisticsPage: statisticsPage,
	}
	go holder.process()
	return holder
}

func NewScriber() *Scriber {
	return NewScriberWithData(&data.Data{ChatStatistics: make(map[int64]*data.ChatStatistics)}, "")
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
			now := time.Now()
			date := now.Format("2006-01-02")
			if chatStatistics.DailyStatistics == nil {
				chatStatistics.DailyStatistics = make(map[string]*data.MessageStatistics)
			}
			dailyStatistics := chatStatistics.DailyStatistics[date]
			if dailyStatistics == nil {
				dailyStatistics = &data.MessageStatistics{}
				chatStatistics.DailyStatistics[date] = dailyStatistics
			}
			dailyStatistics.MessageCounter++
			// Word statistics
			if chatStatistics.DailyWordStatistics == nil {
				chatStatistics.DailyWordStatistics = make(map[string]map[string]int)
			}
			if chatStatistics.DailyWordStatistics[date] == nil {
				chatStatistics.DailyWordStatistics[date] = make(map[string]int)
			}
			words := stem(message.Text)
			for _, word := range words {
				key := normalize(word)
				if key != "" {
					chatStatistics.DailyWordStatistics[date][key] = chatStatistics.DailyWordStatistics[date][key] + 1
				}
			}
			scriber.mutex.Unlock()
		}
	}
}

func stem(text string) []string {
	s := regexp.MustCompile("[^A-Za-z\\p{L}]+").ReplaceAllString(text, " ")
	s = strings.ToLower(s)
	return strings.Fields(s)
}

var prepositions = []string{
	"в", "на", "с", "от", "к", "и", "не", "ну", "он", "так", "там", "то", "что", "чё", "a", "как",
	"the", "of", "and", "to", "in", "is", "it", "et", "by", "from", "or", "but", "has", "that", "are"}

func normalize(word string) string {
	if word == "" {
		return ""
	}
	for _, preposition := range prepositions {
		if preposition == word {
			return ""
		}
	}
	return word
}

func (scriber Scriber) GetStatistics(chatId int64) *data.ChatStatistics {
	return scriber.data.ChatStatistics[chatId]
}

func (scriber Scriber) GetStatisticsPage() string {
	return scriber.statisticsPage
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
	log.Println("usKeys: ", chatStatistics.UsersStatistics)
	usKeys := sortByMessageCounter(chatStatistics.UsersStatistics)
	log.Println("usKeys: ", usKeys)
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

type SortByMessageCounterTuple struct {
	key   string
	value *data.MessageStatistics
}

func sortByMessageCounter(statistics map[string]*data.MessageStatistics) []string {
	tuples := make([]SortByMessageCounterTuple, len(statistics))
	i := 0
	for k := range statistics {
		tuples[i] = SortByMessageCounterTuple{key: k, value: statistics[k]}
		i++
	}
	log.Println("tuples: ", tuples)
	sort.Slice(tuples, func(i, j int) bool {
		return tuples[i].value.MessageCounter > tuples[j].value.MessageCounter
	})
	// take keys
	keys := make([]string, len(tuples))
	for i, tuple := range tuples {
		keys[i] = tuple.key
	}
	return keys
}
