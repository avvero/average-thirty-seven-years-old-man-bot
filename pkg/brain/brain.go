package brain

import (
	"github.com/avvero/the_gamers_guild_bot/internal/knowledge"
	"strings"

	"github.com/avvero/the_gamers_guild_bot/internal/utils"
)

type Brain struct {
	memory       *Memory
	randomFactor bool
}

func NewBrain(memory *Memory, randomFactor bool) *Brain {
	return &Brain{memory: memory, randomFactor: randomFactor}
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
	return with(strings.ToLower(strings.TrimSpace(text))).
		when(truth(brain.randomFactor), random(100)).then(&SenselessPhrasesIntention{}).
		when(truth(brain.randomFactor), random(200), length(5)).then(&HuefyLastWordIntention{}).
		when(truth(brain.randomFactor), random(200), length(14)).then(&HuefyIntention{}).
		when(truth(brain.randomFactor), random(200), length(14)).then(NewKhaleesifyIntention()).
		when(truth(brain.randomFactor), random(10), contains("опять")).say("не опять, а снова").
		when(truth(brain.randomFactor), random(10), contains("купил")).say("А не пиздишь? Аренда это не покупка").
		when(is("gg")).say("gg").
		when(is("нет")).say("пидора ответ").
		when(contains("morrowind", "морровинд", "моровинд")).say("Morrowind - одна из лучших игр эва").
		when(is("er", "ер", "эр")).say("Elden Ring - это величие").
		when(contains("elden ring", "элден ринг", "eлден ринг", "елден ринг", " er ", " ер ", " эр ")).say("Elden Ring - это величие").
		when(contains("spotify", "спотифай")).say("Эти пидоры Антону косарик должны за подписку").
		when(contains("devops", "девопс")).say("Девопсы не нужны").
		when(contains("devops")).say("Девопсы не нужны").
		when(contains("трансформациями")).replace("трансформациями", "оргиями гомогеев").
		when(contains("трансформация")).replace("трансформация", "оргия гомогеев").
		when(contains("трансформацию")).replace("трансформацию", "оргию гомогеев").
		when(contains("трансформации")).replace("трансформации", "оргии гомогеев").
		when(contains("java", "джаба", "джава", "ява")).say("джава-хуява, а я работаю на го").
		when(contains("блокир")).say("пусть себе анус заблокируют").
		when(contains("проблем")).say("у меня есть 5-10 солюшенов этой проблемы").
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

func truth(value bool) func(_ string) bool {
	return func(_ string) bool {
		return value
	}
}

func random(factor int) func(origin string) bool {
	return func(origin string) bool {
		return utils.RandomUpTo(factor) == 0
	}
}

func is(values ...string) func(origin string) bool {
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

func (this *Chain) run() (bool, string) {
	for _, opinion := range this.opinions {
		has, message := opinion.Express(this.text)
		if has && message != "" {
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
