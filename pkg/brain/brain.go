package brain

import (
	"encoding/json"
	"fmt"
	"github.com/avvero/the_gamers_guild_bot/internal/data"
	"github.com/avvero/the_gamers_guild_bot/internal/knowledge"
	"github.com/avvero/the_gamers_guild_bot/internal/openai"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/avvero/the_gamers_guild_bot/internal/utils"
)

type Brain struct {
	memory           *Memory
	randomFactor     bool
	scriber          *statistics.Scriber
	toxicityDetector ToxicityDetector
	openAiClient     *openai.OpenAiClient
}

func NewBrain(randomFactor bool, scriber *statistics.Scriber, toxicityDetector ToxicityDetector, openAiClient *openai.OpenAiClient) *Brain {
	return &Brain{randomFactor: randomFactor, scriber: scriber, toxicityDetector: toxicityDetector, openAiClient: openAiClient}
}

func (brain *Brain) Decision(chatId int64, user string, text string) (respond bool, response string, toxicityScore float64) {
	for _, protector := range []Protector{&Whitelist{}} {
		forbidden, message := protector.Check(chatId, text)
		if forbidden && message != "" {
			return true, message, -1
		} else if forbidden {
			return false, "", -1
		}
	}
	toxicityScore, toxicityDetectionErr := brain.toxicityDetector.ToxicityScore(text)
	if toxicityDetectionErr != nil {
		log.Println("Toxicity check error: " + toxicityDetectionErr.Error())
	}
	respond, response = with(strings.ToLower(strings.TrimSpace(text))).
		// Commands
		when(its("/info")).say("I'm bot").
		when(its("/statistics")).say("не веди себя как магл, используй заклинание: статистика хуистика").
		when(startsWith("/toxicity")).say("не веди себя как магл, используй заклинание: токсик ревиленто").
		when(startsWith("/ai")).say("не веди себя как магл, используй заклинание: интелекто ебанина").
		when(its("статистика хуистика")).say("aхалай махалай, иди сам и посчитай").
		when(its("статистика хуистика")).say(brain.scriber.GetStatisticsPrettyPrint(chatId)).
		when(startsWith("токсик ревиленто")).say(describeToxicity(toxicityScore, toxicityDetectionErr)).
		//when(startsWith("интелекто ебанина"), cost(brain, chatId, user, 1)).say("Репутация: "+strconv.Itoa(brain.scriber.GetUserMessageCount(chatId, user))+". Стоимость навыка: 1. У вас недостаточно репутации для этого этого. Чтобы ее накопить общайтесь или поиграйте с ботом в кости: ролус дайсус.").
		when(startsWith("интелекто ебанина")).say("больше нет, теперь ботян").
		when(startsWith("интеллекто ебанина")).say("больше нет, теперь ботян").
		when(random(50)).then(&OpenApiIntention{brain: brain, model: "gpt-3.5-turbo", text: user + " говорит \"" + text + "\". " +
		"Дай ему краткий ответ в стиле выражений:\n" +
		"Мне похую, я так чувствую\nБлять я заплакал!\n" +
		"Я реально прихуел.\n" +
		"Это было не просто смело, это было пиздец как смело.\n" +
		"Круто! Да это ж круто!\n" +
		"Какое-то величие, какая-то хуйня блять, мне вообще они не интересны\n" +
		"Я уважаю то, что они делают, нооо...\n" +
		"Итс окэй.\n" +
		"Сомнительно, но... Окэй?\n" +
		"Нуууу...\n" +
		"Какая разница?.. Ну вот какая разница?..\n" +
		"Это, конечно, печально. Это печально.\n" +
		"Ни. Хуя. Я. Вот просто: Ни. Хуя. Я.\n" +
		"Ну что это за пиздец такой?\n" +
		"Ну это пиздец какой-то просто! Ну сколько можно?!\n" +
		"Конечно, мы все виноваты в этом пиздеце.."}).
		when(startsWith("поебот")).say("Я теперь Ботян").
		when(its("ботян")).say(user+"?").
		when(startsWith("ботян")).then(&OpenApiIntentionWithError{brain: brain, text: strings.ReplaceAll(text, "ботян", "")}).
		when(startsWith("дайджест")).then(&Digest{brain: brain, chatId: chatId, text: text}).
		when(its("дайджест")).then(&Digest{brain: brain, chatId: chatId, text: text}).
		//when(startsWith("ролус дайсус")).then(&Dice{brain: brain, chatId: chatId, user: user, text: strings.ReplaceAll(text, "ролус дайсус", "")}).
		when(its("мементос ревиленто")).say(brain.scriber.GetNotificationsPrettyPrint(chatId)).
		when(is(brain.randomFactor), random(20), startsWith("мементо")).then(&Notify{brain: brain, chatId: chatId, user: user, text: strings.ReplaceAll(text, "мементо ", ""), action: "Выгнать Вадима"}).
		when(startsWith("мементо")).then(&Notify{brain: brain, chatId: chatId, user: user, text: strings.ReplaceAll(text, "мементо ", "")}).
		when(startsWith("ботян напомни")).then(&Notify{brain: brain, chatId: chatId, user: user, text: strings.ReplaceAll(text, "ботян напомни", "")}).
		//
		when(is(brain.randomFactor), random(30), length(30)).then(&OpenApiIntention{brain: brain, text: user + " говорит \"" + text + "\". Проверь на грамотность, предложи исправление одной ошибки. " +
		"Учти неформальный стиль общения и не будь строг. Если сообщение содержит ошибку, то сообщи об этом. " +
		"Если сообщение не содержит ошибок, то ответь едко как Джимми Карр, но не говори, что ты Джимми Карр."}).
		//when(is(brain.randomFactor), random(500)).then(&DumbledoreScore{brain: brain, chatId: chatId, user: user}).
		when(is(brain.randomFactor), is(toxicityScore >= 0.99)).then(&ToxicReparation{brain: brain, chatId: chatId, user: user}).
		when(is(brain.randomFactor), is(toxicityScore >= 0.99)).then(&OpenApiIntention{brain: brain, model: "gpt-3.5-turbo", text: user + " говорит \"" + text + "\", дай свою оценку его словам и совет согласно тому, что написано в Nonviolent Communication: A Language of Life: Life-Changing Tools for Healthy Relationships. Уложись в пару предложений, не упоминай название книги. А в конце сделай вывод в 1 предложение в стиле  Луи Си Кея, чтобы поставить обидчика на место в грубой форме, не упоминай  Луи Си Кея."}).
		//when(is(brain.randomFactor), is(toxicityScore >= 0.98)).say("на грани щас").
		//when(is(brain.randomFactor), is(toxicityScore >= 0.92)).say("осторожнее").
		when(is(brain.randomFactor), random(100), length(300)).then(&OpenApiIntention{brain: brain, model: "gpt-3.5-turbo", text: "Он говорит \"" + text + "\". , ответить ему, как это бы сделал Джимми Карр."}).
		//when(is(brain.randomFactor), random(50), length(300)).then(&OpenApiIntention{brain: brain, text: "Он говорит \"" + text + "\". , прокомментируй это с издевкой используя морские термины, будто ты пират."}).
		//when(is(brain.randomFactor), random(50), length(100)).then(&OpenApiIntention{brain: brain, text: user + " играет в Dungeons & Dragons, кидает 1d12 кубик. Ты гейммастер, придумай ситуацию, рассчитай сколько на кубике выпало у " + user + ", опиши ситуацию и исход в одно-два предложения, согласно тому сколько на кубике выпало, учитывая то, что сказал " + user + ". " + user + " сказал: \"" + text + "\""}).
		when(is(brain.randomFactor), random(100)).then(&SenselessPhrasesIntention{}).
		when(is(brain.randomFactor), random(100), length(5)).then(&HuefyLastWordIntention{}).
		//when(is(brain.randomFactor), random(300), length(14)).then(&HuefyIntention{}).
		when(is(brain.randomFactor), random(100), length(14)).then(NewKhaleesifyIntention()).
		//when(is(brain.randomFactor), random(1000)).then(&ConfuciusPhrasesIntention{}).
		when(is(brain.randomFactor), random(10), contains("опять")).say("не опять, а снова").
		when(contains("мог быть", "могли быть", "могла быть", "могло быть")).say("словами либерашки").
		when(contains("мнени")).say("Мнение автора может не совпадать с мнением общественности").
		when(its("gg")).say("gg").
		when(its("нет")).then(&NoThenPhraseIntention{}).
		when(contains("morrowind", "морровинд", "моровинд")).say("Morrowind - одна из лучших игр эва").
		when(its("er", "ер", "эр")).say("Elden Ring - это величие").
		when(contains("elden ring", "элден ринг", "eлден ринг", "елден ринг", " er ", " ер ", " эр ")).say("Elden Ring - это величие").
		when(contains("spotify", "спотифай")).say("Эти н'вахи Антону косарик должны за подписку").
		when(is(brain.randomFactor), random(50), contains("devops", "девопс")).say("Девопсы не нужны").
		when(contains("scrum")).say("Скрам - это пережиток").
		when(contains("скрам")).say("Скрам - это пережиток").
		when(contains("предание")).say("Предание - это когда наебали").
		when(is(brain.randomFactor), random(100), contains("трансформациями")).replace("трансформациями", "пертурбациями").
		when(is(brain.randomFactor), random(100), contains("трансформация")).replace("трансформация", "пертурбация").
		when(is(brain.randomFactor), random(100), contains("трансформацию")).replace("трансформацию", "пертурбацию").
		when(is(brain.randomFactor), random(100), contains("трансформации")).replace("трансформации", "пертурбации").
		when(is(brain.randomFactor), random(10), contains("java", "джаба", "джава", "ява")).say("джава-хуява, а я работаю на го").
		//when(contains("блокир")).say("пусть себе анус заблокируют").
		when(is(brain.randomFactor), random(100), contains("проблем")).say("у меня есть 5-10 солюшенов этой проблемы").
		when(contains("mass effect")).say("Шепард умрет").
		when(contains("масс эффект")).say("Шепард умрет").
		//when(contains("новое", "новые", "новая", "новое", "новые", "новье")).say("Точное новое, а не проперженное бу?").
		run()
	return respond, response, toxicityScore
}

func with(text string) *Chain {
	return &Chain{text: text}
}

type Chain struct {
	text     string
	opinions []Opinion
}

type Link struct {
	chain      *Chain
	conditions []func(origin string) bool
}

func (this *Chain) when(conditions ...func(origin string) bool) *Link {
	return &Link{chain: this, conditions: conditions}
}

func (this *Link) then(opinion Opinion) *Chain {
	this.chain.opinions = append(this.chain.opinions, &ConditionOpinion{conditions: this.conditions, opinion: opinion})
	return this.chain
}

func (this *Link) say(text string) *Chain {
	this.chain.opinions = append(this.chain.opinions, &ConditionOpinion{conditions: this.conditions, opinion: &TextOpinion{text: text}})
	return this.chain
}

func (this *Link) replace(from string, to string) *Chain {
	opinion := &TextOpinion{text: strings.Replace(this.chain.text, from, to, -1)}
	this.chain.opinions = append(this.chain.opinions, &ConditionOpinion{conditions: this.conditions, opinion: opinion})
	return this.chain
}

func is(value bool) func(_ string) bool {
	return func(_ string) bool {
		return value
	}
}
func cost(brain *Brain, chatId int64, user string, value int) func(_ string) bool {
	return func(_ string) bool {
		if brain.scriber.GetUserMessageCount(chatId, user) < value {
			return true
		} else {
			brain.scriber.IncreaseUserMessageStatistics(chatId, user, -value)
			return false
		}
	}
}

func random(factor int) func(origin string) bool {
	return func(origin string) bool {
		return utils.RandomUpTo(factor) == 0
	}
}

func its(values ...string) func(origin string) bool {
	return func(origin string) bool {
		for _, value := range values {
			if origin == value {
				return true
			}
			if normalizeRu(origin) == value {
				return true
			}
			if normalizeEn(origin) == value {
				return true
			}
		}
		return false
	}
}

func describeToxicity(score float64, err error) string {
	if err != nil {
		fmt.Printf("Toxicity check error: %s", err)
		return "определить уровень токсичности не удалось, быть может вы - черт, попробуйте позже"
	} else {
		return fmt.Sprintf("уровень токсичности %.0f", score*100) + "%" //TODO
	}
}

func length(size int) func(origin string) bool {
	return func(origin string) bool {
		return len(origin) >= size
	}
}

func contains(values ...string) func(origin string) bool {
	return func(origin string) bool {
		for _, value := range values {
			if strings.Contains(origin, value) {
				return true
			}
			if strings.Contains(normalizeRu(origin), value) {
				return true
			}
			if strings.Contains(normalizeEn(origin), value) {
				return true
			}
		}
		return false
	}
}

func startsWith(value string) func(origin string) bool {
	return func(origin string) bool {
		return strings.HasPrefix(origin, value) && len(strings.TrimSpace(strings.ReplaceAll(origin, value, ""))) > 2
	}
}

func (this *Chain) run() (bool, string) {
	for _, opinion := range this.opinions {
		has, message := opinion.Express(this.text)
		if has && message != "" && message != "null" { // TODO null comes from somewere
			return true, message
		}
	}
	return false, ""
}

func normalizeRu(text string) string {
	result := text
	for k, v := range knowledge.NormalisationMap {
		result = strings.Replace(result, k, v, -1)
	}
	return result
}

func normalizeEn(text string) string {
	result := text
	for k, v := range knowledge.NormalisationMap {
		result = strings.Replace(result, v, k, -1)
	}
	return result
}

type Opinion interface {
	Express(text string) (has bool, response string)
}

type RandomizedOpinion struct {
	factor  int
	opinion Opinion
}

func (this *RandomizedOpinion) Express(text string) (has bool, response string) {
	if utils.RandomUpTo(this.factor) == 0 {
		return this.opinion.Express(text)
	} else {
		return false, ""
	}
}

type ConditionOpinion struct {
	conditions []func(origin string) bool
	opinion    Opinion
}

func (this ConditionOpinion) Express(text string) (has bool, response string) {
	for _, condition := range this.conditions {
		if !condition(text) {
			return false, ""
		}
	}
	return this.opinion.Express(text)
}

type TextOpinion struct {
	text string
}

func (this TextOpinion) Express(text string) (has bool, response string) {
	return true, this.text
}

type ComposeOpinion struct {
	opinions []Opinion
}

func (this ComposeOpinion) Express(text string) (has bool, response string) {
	for _, opinion := range this.opinions {
		has, message := opinion.Express(text)
		if has && message != "" {
			return true, message
		}
	}
	return false, ""
}

type OpenApiIntention struct {
	brain *Brain
	text  string
	model string
}

func (this OpenApiIntention) Express(ignore string) (bool, string) {
	if this.brain.openAiClient == nil {
		return false, ""
	}
	var err error
	var response string
	if this.model != "" {
		err, response = this.brain.openAiClient.CompletionByModel(this.model, this.text)

	} else {
		err, response = this.brain.openAiClient.Completion(this.text)
	}
	if err != nil {
		return false, ""
	} else {
		return true, response
	}
}

type OpenApiIntentionWithError struct {
	brain *Brain
	text  string
}

func (this OpenApiIntentionWithError) Express(ignore string) (has bool, response string) {
	if this.brain.openAiClient == nil {
		return false, ""
	}
	err, response := this.brain.openAiClient.Completion(this.text)
	if err != nil {
		return true, "Ошибка обработки: " + err.Error()
	} else {
		return true, response
	}
}

type Digest struct {
	brain  *Brain
	chatId int64
	text   string
}

func (this Digest) Express(ignore string) (has bool, response string) {
	chatLog := this.brain.scriber.GetChatLog(this.chatId, 200)
	if chatLog == "" {
		return true, "Ничего не происходило"
	}

	err, response := this.brain.openAiClient.CompletionByModel("gpt-4-1106-preview",
		"Представьте, что вы женщина, отвечающая на следующий вопрос. "+
			"Сделай на русском пересказ переписки. Пересказ должен быть информативным, количество слов 200-300. "+
			"Сделай выводы по основным моментам обсуждений. Переписка представлена ниже:\n"+chatLog)
	if err != nil {
		return false, "Ошибка обработки: " + err.Error()
	} else {
		return true, response
	}
}

type DumbledoreScore struct {
	brain  *Brain
	chatId int64
	user   string
}

func (this DumbledoreScore) Express(text string) (has bool, response string) {
	if this.brain.openAiClient == nil {
		return false, ""
	}
	message := "Ученик " + this.user
	score := (utils.RandomUpTo(2) + 1) * 50
	if utils.RandomUpTo(2) == 1 {
		message += " сказал \"" + text + "\"" + " и заработал " + strconv.Itoa(score) + " очков для своего факультета"
		this.brain.scriber.IncreaseUserMessageStatistics(this.chatId, this.user, score)
		//this.brain.scriber.IncreaseHouseScore(this.chatId, userHouse(this.user), score)
	} else {
		message += " сказал ерунду \"" + text + "\"" + " и потерял " + strconv.Itoa(score) + " очков для своего факультета"
		this.brain.scriber.IncreaseUserMessageStatistics(this.chatId, this.user, -score)
		//this.brain.scriber.IncreaseHouseScore(this.chatId, userHouse(this.user), -score)
	}
	message += ", прокомментируй это будто ты профессор Дамблдор с указанием количества очков и названия факультета Хогвартс, где учится " + this.user + " (придумай смешное название факультету в одно слово)"
	err, response := this.brain.openAiClient.Completion(message)
	if err != nil {
		return true, "Давай по новой, " + this.user + ", все хуйня!"
	} else {
		return true, response
	}
}

type Dice struct {
	brain  *Brain
	chatId int64
	user   string
	text   string
}

func (this Dice) Express(ignore string) (has bool, response string) {
	if this.brain.openAiClient == nil {
		return false, ""
	}

	userRoll := utils.RandomUpTo(20) + 1
	gameDescription := "Мы играем в DnD и игрок " + this.user + " бросает кубик 1d20 на ситуацию: '" + this.text + "'. На кубике выпало " + strconv.Itoa(userRoll) + ". Ты гейм-мастер. Придумай, что произойдет и напиши кратко то, что скажет гейм-мастер. Приведи прямую цитату."
	err, response := this.brain.openAiClient.Completion(gameDescription)
	if err != nil {
		return true, "Давай по новой, " + this.user + ", все хуйня!"
	} else {
		return true, response
	}
}

type Notify struct {
	brain  *Brain
	chatId int64
	user   string
	text   string
	action string
}

func (this Notify) Express(ignore string) (has bool, response string) {
	now := time.Now()
	nowString := now.Format("2006-01-02 15:04")
	message := "сейчас " + nowString + ", проанализируй напоминание «" + this.text + "» и представь в json формате с двумя полями: \"action\", \"time\" (в формате yyyy-MM-dd hh:mm)"
	err, response := this.brain.openAiClient.Completion(message)
	if err != nil {
		fmt.Printf("Could not read data: %s\n", err)
		return true, "Давай по новой, " + this.user + ", все хуйня!"
	} else {
		response = strings.ReplaceAll(response, "'", "\"")
		//
		notification := &data.Notification{}
		unmarshalError := json.Unmarshal([]byte(response), notification)
		if unmarshalError != nil {
			fmt.Printf("Could not read data: %s\n", unmarshalError)
			return true, "Не могу разобрать, что ты написал, напиши нормально"
		}
		if notification.Time == "" || notification.Time == nowString {
			return true, "Не могу разобрать, что ты написал, напиши нормально"
		}
		notificationTime, parseTimeError := time.Parse("2006-01-02 15:04", notification.Time)
		if parseTimeError != nil {
			fmt.Printf("Could parse time: %s\n", parseTimeError)
			return true, "Давай по новой, " + this.user + ", все хуйня!"
		}
		if notificationTime.Before(now) {
			return true, "И как ты себе это представляешь, пес?"
		}
		notification.User = this.user
		if -1001733786877 == this.chatId && this.action != "" {
			this.brain.scriber.AddNotification(this.chatId, this.user, this.action, notification.Time)
			return true, "Напомню «" + this.action + "» в " + notification.Time
		} else {
			this.brain.scriber.AddNotification(this.chatId, this.user, notification.Action, notification.Time)
			return true, "Напомню «" + notification.Action + "» в " + notification.Time
		}
	}
}

type ToxicReparation struct {
	brain  *Brain
	chatId int64
	user   string
}

func (this ToxicReparation) Express(ignore string) (has bool, response string) {
	fiveMinutesAgo := time.Now().UTC().Add(-time.Minute * time.Duration(5))
	affectedMessage := strings.Builder{}
	messageBody := "Есть текст: \"Выражаю глубокую озабоченность касательно токсичного поведения " + this.user + ", такое " +
		"поведение могло нанесло моральный ущерб некоторым людям и им будет выплачена компенсация.\". " + userRoleDescription(this.user) +
		"Перескажи так, как это бы сделал директор школы профессор Альбус Дамблдор. Не используй оригинальные слова. Не упоминай название школы."
	err, aiResponse := this.brain.openAiClient.Completion(messageBody)
	if err != nil {
		affectedMessage.WriteString("Выражаю глубокую озабоченность касательно токсичного поведения " + this.user +
			", такое поведение нанесло моральный ущерб некоторым людям. Им будет выплачена компенсация:\n")
	} else {
		affectedMessage.WriteString(aiResponse + "\n")
	}
	userActivity := this.brain.scriber.GetUserActivity(this.chatId)
	usKeys := sortByMessageCounter(userActivity)
	affected := 0
	for _, user := range usKeys {
		if user == this.user {
			continue
		}
		if user == "bot" {
			continue
		}
		lastMessageDateTime := userActivity[user]
		if lastMessageDateTime == "" {
			continue
		}
		dateTime, err := time.ParseInLocation("2006-01-02 15:04:05", lastMessageDateTime, time.Local)
		if err != nil {
			fmt.Printf("Can't parse date %s: %s", lastMessageDateTime, err)
		}
		if dateTime.After(fiveMinutesAgo) {
			affected++
			affectedMessage.WriteString(" - " + user + ": +10\n")
			this.brain.scriber.IncreaseUserMessageStatistics(this.chatId, user, 10)
			this.brain.scriber.IncreaseHouseScore(this.chatId, userHouse(user), 10)
			this.brain.scriber.IncreaseUserMessageStatistics(this.chatId, this.user, -10)
			this.brain.scriber.IncreaseHouseScore(this.chatId, userHouse(this.user), -10)
		}
	}
	if affected == 0 {
		messageBody := "Есть текст: \"Выражаю глубокую озабоченность касательно токсичного поведения " + this.user + ", такое " +
			"поведение могло нанесло моральный ущерб некоторым гражданам, но к счастью ущерб минимальный\". " + userRoleDescription(this.user) +
			"Перескажи так, как это бы сделал директор школы профессор Альбус Дамблдор. Не используй оригинальные слова. Не упоминай название школы. Упомни факультет."
		err, aiResponse := this.brain.openAiClient.Completion(messageBody)
		if err != nil {
			return true, "Выражаю глубокую озабоченность касательно токсичного поведения " + this.user + ", такое поведение могло " +
				"нанесло моральный ущерб некоторым гражданам, \nно к счастью все отделались легким негативом."
		} else {
			return true, aiResponse
		}
	} else {
		return true, affectedMessage.String()
	}
}

func userRoleDescription(user string) string {
	house := userHouse(user)
	if house != "" {
		return "Ученик " + user + " представляет факультет " + house + " школы Чародейства и волшебства Хогвартс."
	} else {
		return "Человек " + user + " -рядовой сквиб без выдающихся качеств и не имеет отношение к школе Чародейства и " +
			"волшебства Хогвартс. Он не является учеником, не обладает магическими навыками и нужно проявить жалость к этому человеку."
	}
}

func userHouse(user string) string {
	switch user {
	case "saintnk":
		return "Когтевран"
	case "svg1007":
		return "Слизерин"
	case "avveroll":
		return "Гриффиндор"
	case "wishpering":
		return "Гриффиндор"
	case "justFirst":
		return "Гриффиндор"
	case "MathJay":
		return "Пуффендуй"
	case "Сторожев Сергей":
		return "Пуффендуй"
	// for test purposes
	case "first":
		return "Пуффендуй"
	case "second":
		return "Гриффиндор"
	default:
		return ""
	}
}

func sortByMessageCounter(activity map[string]string) []string {
	users := make([]string, len(activity))
	i := 0
	for k := range activity {
		users[i] = k
		i++
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i] < users[j]
	})
	return users
}
