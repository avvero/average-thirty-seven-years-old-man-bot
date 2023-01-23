package brain

import (
	"fmt"
	"github.com/avvero/the_gamers_guild_bot/internal/knowledge"
	"github.com/avvero/the_gamers_guild_bot/internal/openai"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"log"
	"strconv"
	"strings"

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
		when(its("статистика хуистика")).say(brain.scriber.GetStatisticsPrettyPrint(chatId)).
		when(startsWith("токсик ревиленто")).say(describeToxicity(toxicityScore, toxicityDetectionErr)).
		when(startsWith("интелекто ебанина"), cost(brain, chatId, user, 10)).say("Репутация: "+strconv.Itoa(brain.scriber.GetUserMessageCount(chatId, user))+". Стоимость навыка: 10. У вас недостаточно репутации для этого этого. Чтобы ее накопить общайтесь или поиграйте с ботом в кости: ролус дайсус.").
		when(startsWith("интелекто ебанина")).then(&OpenApiIntentionWithError{brain: brain, text: strings.ReplaceAll(text, "интелекто ебанина", "")}).
		when(its("ролус дайсус")).then(&Dice{brain: brain, chatId: chatId, user: user}).
		//
		when(is(brain.randomFactor), random(100)).then(&DumbledoreScore{brain: brain, chatId: chatId, user: user}).
		when(is(brain.randomFactor), random(10), is(toxicityScore >= 0.99)).say("токсик").
		when(is(brain.randomFactor), random(10), is(toxicityScore >= 0.98)).say("на грани щас").
		when(is(brain.randomFactor), random(10), is(toxicityScore >= 0.92)).say("осторожнее").
		//when(is(brain.randomFactor), random(100), length(20)).then(&OpenApiIntention{brain: brain, text: "Он говорит \"" + text + "\". , ответить ему в грубой форме с подколкой, как это бы сделал остроумный джентльмен."}).
		when(is(brain.randomFactor), random(200)).then(&SenselessPhrasesIntention{}).
		when(is(brain.randomFactor), random(1000), length(5)).then(&HuefyLastWordIntention{}).
		//when(is(brain.randomFactor), random(300), length(14)).then(&HuefyIntention{}).
		when(is(brain.randomFactor), random(1000), length(14)).then(NewKhaleesifyIntention()).
		when(is(brain.randomFactor), random(1000)).then(&ConfuciusPhrasesIntention{}).
		when(is(brain.randomFactor), random(100), contains("опять")).say("не опять, а снова").
		when(contains("мог быть", "могли быть", "могла быть", "могло быть")).say("словами либерашки").
		when(its("gg")).say("gg").
		when(its("нет")).then(&NoThenPhraseIntention{}).
		when(contains("morrowind", "морровинд", "моровинд")).say("Morrowind - одна из лучших игр эва").
		when(its("er", "ер", "эр")).say("Elden Ring - это величие").
		when(contains("elden ring", "элден ринг", "eлден ринг", "елден ринг", " er ", " ер ", " эр ")).say("Elden Ring - это величие").
		when(contains("spotify", "спотифай")).say("Эти н'вахи Антону косарик должны за подписку").
		when(is(brain.randomFactor), random(50), contains("devops", "девопс")).say("Девопсы не нужны").
		when(contains("scrum")).say("Скрам - это пережиток").
		when(contains("скрам")).say("Скрам - это пережиток").
		when(is(brain.randomFactor), random(100), contains("трансформациями")).replace("трансформациями", "пертурбациями").
		when(is(brain.randomFactor), random(100), contains("трансформация")).replace("трансформация", "пертурбация").
		when(is(brain.randomFactor), random(100), contains("трансформацию")).replace("трансформацию", "пертурбацию").
		when(is(brain.randomFactor), random(100), contains("трансформации")).replace("трансформации", "пертурбации").
		//when(contains("java", "джаба", "джава", "ява")).say("джава-хуява, а я работаю на го").
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
}

func (this OpenApiIntention) Express(ignore string) (has bool, response string) {
	if this.brain.openAiClient == nil {
		return false, ""
	}
	err, response := this.brain.openAiClient.Completion(this.text)
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
	score := (utils.RandomUpTo(2) + 1) * 5
	if utils.RandomUpTo(1) == 0 {
		message += " сказал \"" + text + "\"" + " и заработал " + strconv.Itoa(score) + " очков для своего факультета"
		this.brain.scriber.IncreaseUserMessageStatistics(this.chatId, this.user, score)
	} else {
		message += " сказал ерунду \"" + text + "\"" + " и потерял " + strconv.Itoa(score) + " очков для своего факультета"
		this.brain.scriber.IncreaseUserMessageStatistics(this.chatId, this.user, -score)
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
}

func (this Dice) Express(ignore string) (has bool, response string) {
	if this.brain.openAiClient == nil {
		return false, ""
	}

	cost := 100
	botRoll := utils.RandomUpTo(5) + 1
	userRoll := utils.RandomUpTo(5) + 1
	gameDescription := "Игра в кости. Игрок " + this.user + " выбросил на кубике " + strconv.Itoa(userRoll) + ", Игрок bot выбросил на кубике " + strconv.Itoa(botRoll) + "."
	if userRoll == botRoll {
		gameDescription = gameDescription + ". Результат: Ничья."
	} else if userRoll > botRoll {
		gameDescription = gameDescription + ". Результат: " + this.user + " выиграл " + strconv.Itoa(cost) + " очков."
		this.brain.scriber.IncreaseUserMessageStatistics(this.chatId, this.user, cost)
	} else {
		gameDescription = gameDescription + ". Результат: " + this.user + " проиграл " + strconv.Itoa(cost) + " очков."
		this.brain.scriber.IncreaseUserMessageStatistics(this.chatId, this.user, -cost)
	}

	gameDescription = "прокомментируй игру, будто ты диктор: \"" + gameDescription + "\". Опиши ход игры."
	userMessageCount := this.brain.scriber.GetUserMessageCount(this.chatId, this.user)
	if userMessageCount < 100 {
		gameDescription = gameDescription + "Отрази в своих словах презрение к игроку " + this.user + ", игрок ничто и никто, отрази это."
	} else if userMessageCount < 1000 {
		gameDescription = gameDescription + "Отрази в своих словах глубокую ненависть к игроку " + this.user
	} else if userMessageCount < 3000 {
		gameDescription = gameDescription + "Отрази в своих словах уважение к игроку " + this.user
	} else if userMessageCount > 5000 {
		gameDescription = gameDescription + "Отрази в своих словах восхищение и любовь к игроку " + this.user
	}
	err, response := this.brain.openAiClient.Completion(gameDescription)
	if err != nil {
		return true, "Давай по новой, " + this.user + ", все хуйня!"
	} else {
		return true, response
	}
}
