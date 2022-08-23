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
	text = strings.ToLower(text)

	return with(strings.ToLower(text)).
		when(func(origin string) bool { return brain.randomFactor }, compose([]Opinion{
			random(100, &SenselessPhrasesIntention{}),
			when(len(text) > 5 && !strings.Contains(text, " "), random(200, &HuefyLastWordIntention{})),
			when(len(text) > 5 && !strings.Contains(text, " "), random(200, &HuefyIntention{})),
			when(len(text) > 14, random(200, NewKhaleesifyIntention())),
		})).
		when(is("gg"), say("gg")).
		when(is("нет"), say("пидора ответ")).
		when(contains("morrowind", "морровинд", "моровинд"), say("Morrowind - одна из лучших игр эва")).
		when(oneOf("er", "ер", "эр"), say("Elden Ring - это величие")).
		when(contains("elden ring", "элден ринг", "eлден ринг", "елден ринг", " er ", " ер ", " эр "), say("Elden Ring - это величие")).
		when(contains("купил"), say("А не пиздишь? Аренда это не покупка")).
		when(contains("spotify", "спотифай"), say("Эти пидоры Антону косарик должны за подписку")).
		when(contains("devops", "девопс"), say("Девопсы не нужны")).
		when(contains("devops"), say("Девопсы не нужны")).
		when(contains("трансформациями"), replace(text, "трансформациями", "оргиями гомогеев")).
		when(contains("трансформация"), replace(text, "трансформация", "оргия гомогеев")).
		when(contains("трансформацию"), replace(text, "трансформацию", "оргию гомогеев")).
		when(contains("трансформации"), replace(text, "трансформации", "оргии гомогеев")).
		when(contains("java", "джаба", "джава", "ява"), say("джава-хуява, а я работаю на го")).
		when(contains("блокир"), say("пусть себе анус заблокируют")).
		//random(10, when(has(text, "опять"), say("не опять, а снова"))),
		when(contains("проблем"), say("у меня есть 5-10 солюшенов этой проблемы")).
		run()
}

func with(text string) *Chain {
	return &Chain{text: text}
}

type Chain struct {
	text     string
	opinions []Opinion
}

func (this *Chain) when(condition func(origin string) bool, opinion Opinion) *Chain {
	this.opinions = append(this.opinions, &ConditionOpinion{condition: condition, opinion: opinion})
	return this
}

func is(value string) func(origin string) bool {
	return func(origin string) bool {
		if origin == value {
			return true
		}
		if normalizeRu(origin) == value {
			return true
		}
		if normalizeEn(origin) == value {
			return true
		}
		return false
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

func oneOf(values ...string) func(origin string) bool {
	return func(origin string) bool {
		for _, value := range values {
			if is(value)(origin) {
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

func random(factor int, opinion Opinion) *RandomizedOpinion {
	return &RandomizedOpinion{factor: factor, opinion: opinion}
}

type ConditionOpinion struct {
	condition func(origin string) bool
	opinion   Opinion
}

func (this ConditionOpinion) Express(text string) (has bool, response string) {
	if this.condition(text) {
		return this.opinion.Express(text)
	} else {
		return false, ""
	}
}

func when(condition bool, opinion Opinion) *ConditionOpinion {
	return &ConditionOpinion{condition: func(origin string) bool { return condition }, opinion: opinion}
}

type TextOpinion struct {
	text string
}

func (this TextOpinion) Express(text string) (has bool, response string) {
	return true, this.text
}

func say(text string) *TextOpinion {
	return &TextOpinion{text: text}
}

func replace(text string, from string, to string) *TextOpinion {
	return &TextOpinion{text: strings.Replace(text, from, to, -1)}
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

func compose(opinions []Opinion) *ComposeOpinion {
	return &ComposeOpinion{opinions: opinions}
}
