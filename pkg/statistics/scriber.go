package statistics

import (
	"fmt"
	"github.com/avvero/the_gamers_guild_bot/internal/data"
	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Scriber struct {
	packs          chan *Pack
	mutex          sync.Mutex
	data           *data.Data
	statisticsPage string
}

func NewScriberWithData(data *data.Data, statisticsPage string) *Scriber {
	holder := &Scriber{
		packs:          make(chan *Pack, 100),
		data:           data,
		statisticsPage: statisticsPage,
	}
	go holder.process()
	return holder
}

func NewScriber() *Scriber {
	return NewScriberWithData(&data.Data{ChatStatistics: make(map[int64]*data.ChatStatistics)}, "")
}

type Pack struct {
	message       *telegram.WebhookRequestMessage
	toxicityScore float64
}

func (scriber Scriber) Keep(message *telegram.WebhookRequestMessage, toxicityScore float64) {
	scriber.packs <- &Pack{message: message, toxicityScore: toxicityScore}
}

func (scriber Scriber) process() {
	for {
		select {
		case pack := <-scriber.packs:
			message := pack.message
			scriber.mutex.Lock()
			//
			now := time.Now()
			date := now.Format("2006-01-02")
			// By user
			chatStatistics := scriber.data.ChatStatistics[message.Chat.Id]
			if chatStatistics == nil {
				chatStatistics = &data.ChatStatistics{UsersStatistics: make(map[string]*data.UserMessageStatistics)}
				scriber.data.ChatStatistics[message.Chat.Id] = chatStatistics
			}
			user := scriber.GetUser(message)
			userStatistics := chatStatistics.UsersStatistics[user]
			if userStatistics == nil {
				userStatistics = &data.UserMessageStatistics{}
				chatStatistics.UsersStatistics[user] = userStatistics
			}
			userStatistics.MessageCounter++
			userStatistics.LastMessageDate = date
			userStatistics.LastMessageDateTime = now.Format("2006-01-02 15:04:05")
			userStatistics.ToxicityScore = calculateToxicity(userStatistics.ToxicityScore, pack.toxicityScore)
			// Daily
			if chatStatistics.DailyStatistics == nil {
				chatStatistics.DailyStatistics = make(map[string]*data.DayMessageStatistics)
			}
			dailyStatistics := chatStatistics.DailyStatistics[date]
			if dailyStatistics == nil {
				dailyStatistics = &data.DayMessageStatistics{}
				chatStatistics.DailyStatistics[date] = dailyStatistics
			}
			dailyStatistics.MessageCounter++
			dailyStatistics.ToxicityScore = calculateToxicity(dailyStatistics.ToxicityScore, pack.toxicityScore)
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

func calculateToxicity(prev float64, new float64) float64 {
	if new >= 0.99 {
		return prev + 1
	}
	if new >= 0.98 {
		return prev + 0.5
	}
	if new >= 0.92 {
		return prev + 0.1
	}
	return prev
}

func stem(text string) []string {
	s := regexp.MustCompile("[^A-Za-z\\p{L}]+").ReplaceAllString(text, " ")
	s = strings.ToLower(s)
	return strings.Fields(s)
}

var ignoredWords = []string{
	"в", "на", "с", "от", "к", "и", "не", "ну", "он", "так", "там", "то", "что", "чё", "а", "как", "за", "ни", "у",
	"я", "это", "но", "ты", "все", "по", "же", "из", "бы", "уже", "его", "мой", "про", "меня", "вот", "до", "нет",
	"был", "было", "еще", "ещё", "или", "только", "если", "й", "они", "где", "есть", "мне", "даже", "когда", "да", "их",
	"вы", "ага", "для", "тоже", "прям", "него", "чтобы", "тут",
	"a", "the", "of", "and", "to", "in", "is", "it", "et", "by", "from", "or", "but", "has", "that", "are", "o", "so",
	"for", "on", "as", "an", "not", "no", "t", "s", "http", "https", "www", "com", "ru"}

func normalize(word string) string {
	if word == "" {
		return ""
	}
	if len(word) < 3 {
		return ""
	}
	for _, ignoredWord := range ignoredWords {
		if ignoredWord == word {
			return ""
		}
	}
	return word
}

func (scriber Scriber) GetStatistics(chatId int64) *data.ChatStatistics {
	return scriber.data.ChatStatistics[chatId]
}

func (scriber Scriber) GetUser(message *telegram.WebhookRequestMessage) string {
	user := message.From.Username
	if user == "" {
		user = message.From.LastName + " " + message.From.FirstName
	}
	return user
}

func (scriber Scriber) GetUserStatistics(message *telegram.WebhookRequestMessage) int {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	user := scriber.GetUser(message)
	if scriber.data.ChatStatistics[message.Chat.Id] == nil ||
		scriber.data.ChatStatistics[message.Chat.Id].UsersStatistics[user] == nil {
		return 0
	}
	return scriber.data.ChatStatistics[message.Chat.Id].UsersStatistics[user].MessageCounter
}

func (scriber Scriber) GetUserMessageCount(chatId int64, user string) int {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	if scriber.data.ChatStatistics[chatId] == nil ||
		scriber.data.ChatStatistics[chatId].UsersStatistics[user] == nil {
		return 0
	}
	return scriber.data.ChatStatistics[chatId].UsersStatistics[user].MessageCounter
}

func (scriber Scriber) SetUserStatistics(message *telegram.WebhookRequestMessage, messageCounter int) {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	user := scriber.GetUser(message)
	if scriber.data.ChatStatistics[message.Chat.Id] == nil {
		scriber.data.ChatStatistics[message.Chat.Id] = &data.ChatStatistics{UsersStatistics: make(map[string]*data.UserMessageStatistics)}
	}
	if scriber.data.ChatStatistics[message.Chat.Id].UsersStatistics[user] == nil {
		scriber.data.ChatStatistics[message.Chat.Id].UsersStatistics[user] = &data.UserMessageStatistics{}
	}
	scriber.data.ChatStatistics[message.Chat.Id].UsersStatistics[user].MessageCounter = messageCounter
}

func (scriber Scriber) IncreaseUserMessageStatistics(chatId int64, user string, value int) {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	if scriber.data.ChatStatistics[chatId] == nil ||
		scriber.data.ChatStatistics[chatId].UsersStatistics[user] == nil {
		return
	}
	scriber.data.ChatStatistics[chatId].UsersStatistics[user].MessageCounter += value
}

func (scriber Scriber) IncreaseHouseScore(chatId int64, house string, value int) {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	if house == "" {
		return
	}
	if scriber.data.ChatStatistics[chatId] == nil {
		return
	}
	if scriber.data.ChatStatistics[chatId].HouseStatistics == nil {
		scriber.data.ChatStatistics[chatId].HouseStatistics = make(map[string]*data.HouseStatistics)
	}
	if scriber.data.ChatStatistics[chatId].HouseStatistics[house] == nil {
		scriber.data.ChatStatistics[chatId].HouseStatistics[house] = &data.HouseStatistics{}
	}
	scriber.data.ChatStatistics[chatId].HouseStatistics[house].Score += value
}

func (scriber Scriber) GetUserTension(chatId int64, user string) int {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	if scriber.data.ChatStatistics[chatId] == nil ||
		scriber.data.ChatStatistics[chatId].UsersStatistics[user] == nil {
		return 0
	}
	return scriber.data.ChatStatistics[chatId].UsersStatistics[user].Tension
}

func (scriber Scriber) SetUserTension(chatId int64, user string, tension int) {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	if scriber.data.ChatStatistics[chatId] == nil {
		scriber.data.ChatStatistics[chatId] = &data.ChatStatistics{UsersStatistics: make(map[string]*data.UserMessageStatistics)}
	}
	if scriber.data.ChatStatistics[chatId].UsersStatistics[user] == nil {
		scriber.data.ChatStatistics[chatId].UsersStatistics[user] = &data.UserMessageStatistics{}
	}
	scriber.data.ChatStatistics[chatId].UsersStatistics[user].Tension = tension
}

func (scriber Scriber) SetUserLastMessageDate(chatId int64, user string, date string) {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	if scriber.data.ChatStatistics[chatId] == nil ||
		scriber.data.ChatStatistics[chatId].UsersStatistics[user] == nil {
		return
	}
	scriber.data.ChatStatistics[chatId].UsersStatistics[user].LastMessageDate = date
}

func (scriber Scriber) SetUserLastMessageDateTime(chatId int64, user string, dateTime string) {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	if scriber.data.ChatStatistics[chatId] == nil {
		scriber.data.ChatStatistics[chatId] = &data.ChatStatistics{UsersStatistics: make(map[string]*data.UserMessageStatistics)}
	}
	if scriber.data.ChatStatistics[chatId].UsersStatistics[user] == nil {
		scriber.data.ChatStatistics[chatId].UsersStatistics[user] = &data.UserMessageStatistics{}
	}
	scriber.data.ChatStatistics[chatId].UsersStatistics[user].LastMessageDateTime = dateTime
}

func (scriber Scriber) GetStatisticsPage() string {
	return scriber.statisticsPage
}

// TODO move to external object
func (scriber Scriber) GetStatisticsPrettyPrint(chatId int64) string {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	if scriber.data.ChatStatistics == nil ||
		scriber.data.ChatStatistics[chatId] == nil ||
		scriber.data.ChatStatistics[chatId].UsersStatistics == nil ||
		scriber.data.ChatStatistics[chatId].DailyStatistics == nil {
		return ""
	}
	chatStatistics := scriber.data.ChatStatistics[chatId]

	sb := strings.Builder{}
	sb.WriteString("Top 10 users:\n")
	archivedCounter := 0
	archivedMessage := strings.Builder{}
	usKeys := sortByMessageCounter(chatStatistics.UsersStatistics)
	topUsersCounter := 0
	for i := 0; i < len(usKeys) && topUsersCounter <= 10; i++ {
		messages := strconv.Itoa(chatStatistics.UsersStatistics[usKeys[i]].MessageCounter)
		toxicity := fmt.Sprintf("%.2f", chatStatistics.UsersStatistics[usKeys[i]].ToxicityScore)

		if chatStatistics.UsersStatistics[usKeys[i]].LastMessageDate == "" {
			archivedCounter++
			archivedMessage.WriteString(" - " + usKeys[i] + ": " + messages + " (t: " + toxicity + "), last message date: unknown" + "\n")
		} else {
			lastMessageDate, _ := time.Parse("2006-01-02", chatStatistics.UsersStatistics[usKeys[i]].LastMessageDate)
			if time.Now().AddDate(0, 0, -7).After(lastMessageDate) {
				archivedCounter++
				archivedMessage.WriteString(" - " + usKeys[i] + ": " + messages + " (t: " + toxicity + "), last message date: " +
					chatStatistics.UsersStatistics[usKeys[i]].LastMessageDate + "\n")
			} else {
				topUsersCounter++
				sb.WriteString(" - " + usKeys[i] + ": " + messages + " (t: " + toxicity + ")" + "\n")
			}
		}
	}
	sb.WriteString("\n")
	sb.WriteString("Last 10 days:\n")
	dsKeys := sortedKeys(chatStatistics.DailyStatistics)
	start := 0
	if len(dsKeys) > 10 {
		start = len(dsKeys) - 10
	}
	for i := start; i < len(dsKeys); i++ {
		messages := strconv.Itoa(chatStatistics.DailyStatistics[dsKeys[i]].MessageCounter)
		toxicity := fmt.Sprintf("%.2f", chatStatistics.DailyStatistics[dsKeys[i]].ToxicityScore)
		sb.WriteString(" - " + dsKeys[i] + ": " + messages + " (t: " + toxicity + ")" + "\n")
	}

	showInfuriating := false
	for i := 0; i < len(usKeys); i++ {
		if chatStatistics.UsersStatistics[usKeys[i]].Tension > 0 {
			showInfuriating = true
			break
		}
	}
	if showInfuriating {
		sb.WriteString("\n")
		sb.WriteString("Top 10 infuriating persons:\n")

		for i := 0; i < len(usKeys); i++ {
			if chatStatistics.UsersStatistics[usKeys[i]].Tension == 0 {
				continue
			}
			tension := strconv.Itoa(chatStatistics.UsersStatistics[usKeys[i]].Tension)
			sb.WriteString(" - " + usKeys[i] + ": tension = " + tension + "\n")
		}
	}
	// Houses
	houseKeys := sortByScore(chatStatistics.HouseStatistics)
	if len(houseKeys) > 0 {
		sb.WriteString("\n")
		sb.WriteString("Houses scores:\n")
		for i := 0; i < len(houseKeys); i++ {
			if chatStatistics.HouseStatistics[houseKeys[i]].Score == 0 {
				continue
			}
			score := strconv.Itoa(chatStatistics.HouseStatistics[houseKeys[i]].Score)
			sb.WriteString(" - " + houseKeys[i] + ": " + score + "\n")
		}
	}

	// Archived
	if archivedCounter > 0 {
		sb.WriteString("\n")
		sb.WriteString("Between-Morrowind-And-Skyrim: " + strconv.Itoa(archivedCounter) + "\n")
		//sb.WriteString(archivedMessage.String())
	}
	sb.WriteString("\n")
	sb.WriteString("To get more information visit: " + scriber.GetStatisticsPage() + "?id=" + strconv.FormatInt(chatId, 10))
	return sb.String()
}

func (scriber Scriber) GetNotificationsPrettyPrint(chatId int64) string {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	notifications := scriber.GetNotifications(chatId)
	if notifications != nil && len(notifications) > 0 {
		sb := strings.Builder{}
		sb.WriteString("Notifications:\n")

		for key, value := range notifications {
			sb.WriteString(" - " + key + ": " + value.Action + "\n")
		}
		return sb.String()
	} else {
		return "Ничего нет, просто ничего, просто 0, ни ху я"
	}
}

func (scriber Scriber) GetNotifications(chatId int64) map[string]*data.Notification {
	if scriber.data.ChatStatistics[chatId] == nil {
		return nil
	}
	return scriber.data.ChatStatistics[chatId].Notifications
}

func (scriber Scriber) GetChatStatistics() map[int64]*data.ChatStatistics {
	if scriber.data == nil {
		return nil
	}
	return scriber.data.ChatStatistics
}

func (scriber Scriber) AddNotification(chatId int64, user string, action string, time string) {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	if scriber.data.ChatStatistics[chatId] == nil {
		scriber.data.ChatStatistics[chatId] = &data.ChatStatistics{UsersStatistics: make(map[string]*data.UserMessageStatistics)}
	}
	if scriber.data.ChatStatistics[chatId].Notifications == nil {
		scriber.data.ChatStatistics[chatId].Notifications = make(map[string]*data.Notification)
	}
	notification := &data.Notification{User: user, Action: action, Time: time}
	scriber.data.ChatStatistics[chatId].Notifications[notification.Time] = notification
}

func (scriber Scriber) RemoveNotification(chatId int64, time string) {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	if scriber.data.ChatStatistics[chatId] == nil || scriber.data.ChatStatistics[chatId].Notifications == nil {
		return
	}
	delete(scriber.data.ChatStatistics[chatId].Notifications, time)
}

func (scriber Scriber) GetUserActivity(chatId int64) map[string]string {
	scriber.mutex.Lock()
	defer scriber.mutex.Unlock()

	result := make(map[string]string)
	if scriber.GetStatistics(chatId) != nil {
		for user, statistics := range scriber.GetStatistics(chatId).UsersStatistics {
			result[user] = statistics.LastMessageDateTime
		}
	}
	return result
}

func sortedKeys(m map[string]*data.DayMessageStatistics) []string {
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
	value *data.UserMessageStatistics
}

func sortByMessageCounter(statistics map[string]*data.UserMessageStatistics) []string {
	tuples := make([]SortByMessageCounterTuple, len(statistics))
	i := 0
	for k := range statistics {
		tuples[i] = SortByMessageCounterTuple{key: k, value: statistics[k]}
		i++
	}
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

type SortByScoreTuple struct {
	key   string
	value *data.HouseStatistics
}

func sortByScore(statistics map[string]*data.HouseStatistics) []string {
	tuples := make([]SortByScoreTuple, len(statistics))
	i := 0
	for k := range statistics {
		tuples[i] = SortByScoreTuple{key: k, value: statistics[k]}
		i++
	}
	sort.Slice(tuples, func(i, j int) bool {
		return tuples[i].value.Score > tuples[j].value.Score
	})
	// take keys
	keys := make([]string, len(tuples))
	for i, tuple := range tuples {
		keys[i] = tuple.key
	}
	return keys
}
