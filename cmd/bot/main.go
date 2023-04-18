package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/avvero/the_gamers_guild_bot/internal/data"
	"github.com/avvero/the_gamers_guild_bot/internal/huggingface"
	"github.com/avvero/the_gamers_guild_bot/internal/openai"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/avvero/the_gamers_guild_bot/internal/telegram"
	"github.com/avvero/the_gamers_guild_bot/internal/utils"
	"github.com/avvero/the_gamers_guild_bot/pkg/brain"

	_ "github.com/motemen/go-loghttp/global" // Just this line!
)

var (
	httpPort             = flag.String("httpPort", "8080", "http server port")
	token                = flag.String("token", "PROVIDE", "bot token")
	jsonBinMasterKey     = flag.String("jsonBinMasterKey", "PROVIDE", "jsonBinMasterKey")
	huggingfaceAccessKey = flag.String("huggingfaceAccessKey", "PROVIDE", "huggingfaceAccessKey")
	statisticsPage       = flag.String("statistics-page", "PROVIDE", "statistics-page")
	openApiKey           = flag.String("open-api-key", "PROVIDE", "open-api-key")
)

func main() {
	gracefullShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefullShutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	flag.Parse()

	tokenEnv, found := os.LookupEnv("token")
	if found {
		token = &tokenEnv
	}
	jsonBinMasterKeyEnv, found := os.LookupEnv("json-bin-master-key")
	if found {
		jsonBinMasterKey = &jsonBinMasterKeyEnv
	}
	statisticsPageEnv, found := os.LookupEnv("statistics-page")
	if found {
		statisticsPage = &statisticsPageEnv
	}
	//Data
	jsonBinClient := data.NewJsonBinApiClient(*jsonBinMasterKey)
	data, err := jsonBinClient.Read()
	if err != nil {
		fmt.Printf("Could not read data: %s\n", err)
		panic(err)
	}
	ticker := time.NewTicker(1 * time.Hour)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case t := <-ticker.C:
				fmt.Println("Write data to bin", t)
				//sendMessage(245851441, 0, "Write data to bin")
				err := jsonBinClient.Write(data)
				if err != nil {
					sendMessage(245851441, 0, "Write data to bin erro: "+err.Error())
				}
			}
		}
	}()
	scriber := statistics.NewScriberWithData(data, *statisticsPage)
	// Toxicity detector
	huggingfaceAccessKeyEnv, found := os.LookupEnv("huggingface-access-key")
	if found {
		huggingfaceAccessKey = &huggingfaceAccessKeyEnv
	}
	url := "https://api-inference.huggingface.co/models/apanc/russian-inappropriate-messages"
	huggingFaceApiClient := huggingface.NewApiClient(url, huggingfaceAccessKeyEnv)
	toxicityDetector := brain.NewToxicityDetector(huggingFaceApiClient)
	// open api
	openApiKeyEnv, found := os.LookupEnv("open-api-key")
	if found {
		openApiKey = &openApiKeyEnv
	}
	openApiClient := openai.NewApiClient("https://api.openai.com/v1/chat/completions", openApiKeyEnv)
	//
	brain := brain.NewBrain(true, scriber, toxicityDetector, &openApiClient)

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{\"name\": \"Average Thirty-Seven Years Old Man (bot)\", \"version\": \"1.4\"}")
	})

	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/statistics", func(w http.ResponseWriter, r *http.Request) {
		chatIdString := r.URL.Query().Get("id")
		chatId, _ := strconv.ParseInt(chatIdString, 10, 64)

		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, utils.PrintJson(scriber.GetStatistics(chatId)))
	})

	http.HandleFunc("/discord", func(w http.ResponseWriter, r *http.Request) {
		guildId := r.URL.Query().Get("guildId")
		channelId := r.URL.Query().Get("channelId")
		http.Redirect(w, r, fmt.Sprintf("discord://-/channels/%s/%s", guildId, channelId), 301)
	})

	http.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("could not read body: %s\n", err)
			return
		}
		fmt.Println("Telegram Message body from " + string(body))
		webhookRequest := &telegram.WebhookRequest{}
		json.Unmarshal(body, webhookRequest)
		if webhookRequest == nil || webhookRequest.Message == nil {
			fmt.Printf("could not unmarshal body\n")
			return
		}
		fmt.Println("Message from " + strconv.FormatInt(webhookRequest.Message.Chat.Id, 10) + " " +
			webhookRequest.Message.Chat.Title + ": " + webhookRequest.Message.Text)
		if webhookRequest.Message.ForwardFromMessageId != 0 {
			fmt.Println("Skip forward message")
			return
		}
		if webhookRequest.Message.NewChatParticipant != nil && "HatefulVadimBot" == webhookRequest.Message.NewChatParticipant.Username {
			go func() {
				fmt.Println("HatefulVadimBot is the new participant")
				sendMessage(webhookRequest.Message.Chat.Id, 0, "Так и кто это?")
				sendChatAction(webhookRequest.Message.Chat.Id, "typing")
				err, aiResponse := openApiClient.Completion("Есть текст: \"В нашем сообществе есть невыносимо неприятная " +
					"личность, я не могу больше терпеть его присутствие и просто обязан его выгнать\". Перескажи так, " +
					"как это бы сделал директор школы профессор Альбус Дамблдор. Не используй оригинальные слова.")
				if err != nil {
					sendMessage(webhookRequest.Message.Chat.Id, 0, "Не знаю как сказать это красиво, но давайте его нахер выгонем?")
				} else {
					sendMessage(webhookRequest.Message.Chat.Id, 0, aiResponse)
				}
				time.Sleep(time.Duration(utils.RandomUpTo(2)) * time.Second)
				banChatMember(webhookRequest.Message.Chat.Id, webhookRequest.Message.NewChatParticipant.Id)
			}()
			return
		}
		user := scriber.GetUser(webhookRequest.Message)
		go func() {
			respond, response, toxicityScore := brain.Decision(webhookRequest.Message.Chat.Id, user, webhookRequest.Message.Text)
			scriber.Keep(webhookRequest.Message, toxicityScore)
			fmt.Println("Message to " + strconv.FormatInt(webhookRequest.Message.Chat.Id, 10) + " " +
				webhookRequest.Message.Chat.Title + ": " + strconv.FormatBool(respond) + ": " + response)
			if respond {
				//time.Sleep(time.Duration(utils.RandomUpTo(15)) * time.Second)
				sendMessage(webhookRequest.Message.Chat.Id, webhookRequest.Message.MessageId, response)
				// wrap
				botMessage := &telegram.WebhookRequestMessage{
					MessageId: 0,
					From:      &telegram.WebhookRequestMessageSender{Username: "bot"},
					Chat:      webhookRequest.Message.Chat,
					Text:      response,
				}
				scriber.Keep(botMessage, 0)
			}
		}()

		//
		if scriber.GetUserStatistics(webhookRequest.Message) >= 10000 {
			user := scriber.GetUser(webhookRequest.Message)
			sendMessage(webhookRequest.Message.Chat.Id, 0, "Пользователь "+user+" достиг величия и начинает свой новый цикл восхождения. Путем унижения он пройдет вновь от низших и бесполезных к великим и бесценным.")
			scriber.SetUserStatistics(webhookRequest.Message, 0)
			jsonBinClient.Write(data)
		}
	})
	// Scheduler
	standupScheduler := gocron.NewScheduler(time.UTC)
	standupScheduler.Cron("13 7 * * 1-5").Do(func() {
		sendMessage(-1001733786877, 0, "@avveroll, @wishpering, @justFirst пиздуйте на стэндап")
	})
	standupScheduler.StartAsync()
	// Notifications
	notificationsTicker := time.NewTicker(1 * time.Minute)
	notificationDone := make(chan bool)
	go func() {
		for {
			select {
			case <-notificationDone:
				ticker.Stop()
				return
			case t := <-notificationsTicker.C:
				fmt.Println("Check notifications, time: ", t)
				for chatId, _ := range scriber.GetChatStatistics() {
					notifications := scriber.GetNotifications(chatId)
					if notifications != nil {
						for _, notification := range notifications {
							notificationTime, parseTimeError := time.Parse("2006-01-02 15:04", notification.Time)
							if parseTimeError != nil {
								fmt.Printf("Could parse time: %s\n", parseTimeError)
								continue
							}
							if notificationTime.Before(t) {
								fmt.Printf("Time is passed for " + notification.Time + ": " + notification.Action)
								sendMessage(chatId, 0, "Напоминаю по просьбе @"+notification.User+": "+notification.Action)
								scriber.RemoveNotification(chatId, notification.Time)
							}
						}
					}
				}
			}
		}
	}()

	// discord api
	discordBoyKeyEnv, found := os.LookupEnv("discord-bot-key")
	if !found {
		fmt.Println("Can't find discord bot api key")
		os.Exit(0)
	}
	discord, err := discordgo.New("Bot " + discordBoyKeyEnv)
	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	discord.AddHandler(messageCreate)
	discord.AddHandler(presenceUpdate(&openApiClient))
	discord.AddHandler(PresencesReplace)
	discord.AddHandler(VoiceStateUpdate(*statisticsPage)) //
	discord.Identify.Intents = discordgo.IntentsAll

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		panic(err)
	}
	defer discord.Close()

	log.Println("Http server started on port " + *httpPort)
	sendMessage(245851441, 0, "Bot is started, version 1.7")
	err, aiResponse := openApiClient.Completion("Придумай остроумное приветствие")
	if err != nil {
		sendMessage(245851441, 0, "Ошибка AI: "+err.Error())
	} else {
		sendMessage(245851441, 0, aiResponse)
	}
	http.ListenAndServe(":"+*httpPort, nil)
	<-gracefullShutdown
	jsonBinClient.Write(data)
	sendMessage(245851441, 0, "Bot is stopped, version 1.7")
}

func messageCreate(s *discordgo.Session, event *discordgo.MessageCreate) {
	payload, _ := json.Marshal(event)
	// Set the playing status.
	fmt.Printf("Incoming event %s\n", string(payload))
	fmt.Printf("Discord message: %s: %s\n", event.Author.Username, event.Content)
	sendMessage(245851441, 0, fmt.Sprintf("Discord message: %s: %s\n", event.Author.Username, event.Content))
}

func presenceUpdate(openAiClient *openai.OpenAiClient) func(s *discordgo.Session, event *discordgo.PresenceUpdate) {
	activityMap := make(map[string]string)
	return func(s *discordgo.Session, event *discordgo.PresenceUpdate) {
		payload, _ := json.Marshal(event)
		fmt.Printf("Incoming event %s\n", string(payload))
		// Set the playing status.
		userId := event.Presence.User.ID
		user, _ := s.User(userId)
		if len(event.Presence.Activities) > 0 {
			game := event.Presence.Activities[0].Name
			fmt.Println("Discord activity start: ", user.Username, event.Presence.Status, game)
			if event.Presence.Activities[0].Type == discordgo.ActivityTypeGame && activityMap[userId] != game {
				message := fmt.Sprintf("Есть новость: %s начал играть в %s. Расскажи об этом коротко в едкой форме.", user.Username, game)
				err, aiResponse := openAiClient.Completion(message)
				if err != nil {
					sendMessage(245851441, 0, "Ошибка AI: "+err.Error())
				} else {
					sendMessage(-1001733786877, 0, aiResponse)
					activityMap[userId] = game
				}
			}
		} else {
			fmt.Println("Discord activity stop: ", user.Username)
			if activityMap[userId] != "" {
				message := fmt.Sprintf("Есть новость: %s закончил играть в %s. Расскажи об этом коротко в едкой форме.",
					user.Username, activityMap[userId])
				err, aiResponse := openAiClient.Completion(message)
				if err != nil {
					sendMessage(245851441, 0, "Ошибка AI: "+err.Error())
				} else {
					sendMessage(-1001733786877, 0, aiResponse)
				}
				activityMap[userId] = ""
			}
			//sendMessage(245851441, 0, fmt.Sprint("Discord activity stop: ", user.Username))
		}
	}
}

func PresencesReplace(s *discordgo.Session, presencies []*discordgo.Presence) {
	payload, _ := json.Marshal(presencies)
	fmt.Printf("presencies %#v\n", string(payload))
}

func VoiceStateUpdate(domain string) func(s *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	return func(s *discordgo.Session, event *discordgo.VoiceStateUpdate) {
		payload, _ := json.Marshal(event)
		if event.ChannelID != "" {
			user, _ := s.User(event.UserID)
			channel, _ := s.Channel(event.ChannelID)
			fmt.Printf("VoiceStateUpdate %s\n", string(payload))
			//
			sendMessage(-1001733786877, 0, fmt.Sprintf("%s зашел в голосовой канал discord сервера: [%s](%s)", user.Username,
				channel.Name, fmt.Sprintf("%s/discord?guildId=%s&channelId=%s", domain, channel.GuildID, channel.ID)))
		}
	}
}

func sendMessage(chatId int64, receivedMessageId int64, message string) {
	replyToMessageId := ""
	if receivedMessageId != 0 {
		replyToMessageId = strconv.FormatInt(receivedMessageId, 10)
	}
	requestBody, marshalError := json.Marshal(map[string]string{
		"reply_to_message_id": replyToMessageId,
		"chat_id":             strconv.FormatInt(chatId, 10),
		"text":                message,
		"parse_mode":          "markdown",
	})
	if marshalError != nil {
		fmt.Printf("could not marshal body: %s\n", marshalError)
	}
	client := http.Client{Timeout: 5 * time.Second}
	url := "https://api.telegram.org/bot" + *token + "/sendMessage"
	fmt.Printf("Request to: %s, message: %s\n", url, message)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	_, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return
	}
}

func sendSticker(chatId int64, receivedMessageId int64, fileId string) {
	replyToMessageId := ""
	if receivedMessageId != 0 {
		replyToMessageId = strconv.FormatInt(receivedMessageId, 10)
	}
	requestBody, marshalError := json.Marshal(map[string]string{
		"reply_to_message_id": replyToMessageId,
		"chat_id":             strconv.FormatInt(chatId, 10),
		"sticker":             fileId,
	})
	if marshalError != nil {
		fmt.Printf("could not marshal body: %s\n", marshalError)
	}
	client := http.Client{Timeout: 5 * time.Second}
	url := "https://api.telegram.org/bot" + *token + "/sendSticker"
	fmt.Printf("Request to: %s, sticker: %s\n", url, fileId)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	_, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return
	}
}

func banChatMember(chatId int64, userId int64) {
	requestBody, marshalError := json.Marshal(map[string]string{
		"chat_id": strconv.FormatInt(chatId, 10),
		"user_id": strconv.FormatInt(userId, 10),
	})
	if marshalError != nil {
		fmt.Printf("could not marshal body: %s\n", marshalError)
	}
	client := http.Client{Timeout: 5 * time.Second}
	url := "https://api.telegram.org/bot" + *token + "/banChatMember"
	fmt.Printf("Request to: %s, user: %s\n", url, strconv.FormatInt(userId, 10))
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	_, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return
	}
}

func sendChatAction(chatId int64, action string) {
	requestBody, marshalError := json.Marshal(map[string]string{
		"chat_id": strconv.FormatInt(chatId, 10),
		"action":  action,
	})
	if marshalError != nil {
		fmt.Printf("could not marshal body: %s\n", marshalError)
	}
	client := http.Client{Timeout: 5 * time.Second}
	url := "https://api.telegram.org/bot" + *token + "/sendChatAction"
	fmt.Printf("Request to: %s, action: %s\n", url, action)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	_, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
		return
	}
}
