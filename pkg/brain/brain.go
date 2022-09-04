package brain

import (
	"fmt"
	"github.com/avvero/the_gamers_guild_bot/internal/knowledge"
	"github.com/avvero/the_gamers_guild_bot/pkg/statistics"
	"strings"

	"github.com/avvero/the_gamers_guild_bot/internal/utils"
)

type Brain struct {
	memory           *Memory
	randomFactor     bool
	scriber          *statistics.Scriber
	toxicityDetector ToxicityDetector
}

func NewBrain(randomFactor bool, scriber *statistics.Scriber, toxicityDetector ToxicityDetector) *Brain {
	return &Brain{randomFactor: randomFactor, scriber: scriber, toxicityDetector: toxicityDetector}
}

func (brain *Brain) Decision(chatId int64, text string) (respond bool, response string) {
	for _, protector := range []Protector{&Whitelist{}, &Censor{}} {
		forbidden, message := protector.Check(chatId, text)
		if forbidden && message != "" {
			return true, message
		} else if forbidden {
			return false, ""
		}
	}
	toxicityScore, toxicityDetectionErr := brain.toxicityDetector.ToxicityScore(text)
	if toxicityDetectionErr != nil {
		fmt.Printf("Toxicity check error: %s", toxicityDetectionErr)
	}
	return with(strings.ToLower(strings.TrimSpace(text))).
		// Commands
		when(its("/info")).say("I'm bot").
		when(its("/statistics")).say(utils.PrintJson(brain.scriber.GetStatistics(chatId))).
		when(startsWith("/toxicity")).say(describeToxicity(toxicityScore, toxicityDetectionErr)).
		//
		when(is(toxicityScore >= 0.98)).say("токсик").
		when(is(toxicityScore >= 0.90)).say("на грани щас").
		when(is(toxicityScore >= 0.8)).say("осторожнее").
		when(is(brain.randomFactor), random(200)).then(&SenselessPhrasesIntention{}).
		when(is(brain.randomFactor), random(300), length(5)).then(&HuefyLastWordIntention{}).
		when(is(brain.randomFactor), random(300), length(14)).then(&HuefyIntention{}).
		when(is(brain.randomFactor), random(300), length(14)).then(NewKhaleesifyIntention()).
		when(is(brain.randomFactor), random(500)).then(&ConfuciusPhrasesIntention{}).
		when(is(brain.randomFactor), random(100), contains("опять")).say("не опять, а снова").
		when(is(brain.randomFactor), random(100), contains("купил")).say("А не пиздишь? Аренда это не покупка").
		when(its("gg")).say("gg").
		when(its("нет")).say("пидора ответ").
		when(contains("morrowind", "морровинд", "моровинд")).say("Morrowind - одна из лучших игр эва").
		when(its("er", "ер", "эр")).say("Elden Ring - это величие").
		when(contains("elden ring", "элден ринг", "eлден ринг", "елден ринг", " er ", " ер ", " эр ")).say("Elden Ring - это величие").
		when(contains("spotify", "спотифай")).say("Эти пидоры Антону косарик должны за подписку").
		when(contains("devops", "девопс")).say("Девопсы не нужны").
		when(contains("devops")).say("Девопсы не нужны").
		when(contains("scrum")).say("Скрам - это пережиток").
		when(contains("скрам")).say("Скрам - это пережиток").
		when(contains("трансформациями")).replace("трансформациями", "оргиями гомогеев").
		when(contains("трансформация")).replace("трансформация", "оргия гомогеев").
		when(contains("трансформацию")).replace("трансформацию", "оргию гомогеев").
		when(contains("трансформации")).replace("трансформации", "оргии гомогеев").
		when(contains("java", "джаба", "джава", "ява")).say("джава-хуява, а я работаю на го").
		when(contains("блокир")).say("пусть себе анус заблокируют").
		when(contains("проблем")).say("у меня есть 5-10 солюшенов этой проблемы").
		when(contains("mass effect")).say("Шепард умрет").
		when(contains("масс эффект")).say("Шепард умрет").
		run()
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
