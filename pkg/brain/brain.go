package brain

import (
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
	for _, opinion := range []Opinion{
		when(brain.randomFactor, compose([]Opinion{
			random(100, &SenselessPhrasesIntention{}),
			when(len(text) > 5 && !strings.Contains(text, " "), random(200, &HuefyLastWordIntention{})),
			when(len(text) > 5 && !strings.Contains(text, " "), random(200, &HuefyIntention{})),
			when(len(text) > 14, random(200, NewKhaleesifyIntention())),
		})),
		when(text == "gg", say("gg")),
		when(brain.normalizeRu(text) == "нет", say("пидора ответ")),
		when(has(text, []string{"morrowind", "морровинд", "моровинд"}), say("Morrowind - одна из лучших игр эва")),
		when(isOf(text, []string{"er", "ер", "эр"}), say("Elden Ring - это величие")),
		when(has(text, []string{"elden ring", "элден ринг", "eлден ринг", "елден ринг", " er ", " ер ", " эр "}), say("Elden Ring - это величие")),
		when(has(text, []string{"купил"}), say("А не пиздишь? Аренда это не покупка")),
		when(has(text, []string{"spotify", "спотифай"}), say("Эти пидоры Антону косарик должны за подписку")),
		when(has(brain.normalizeEn(text), []string{"devops"}), say("Девопсы не нужны")),
		when(has(brain.normalizeRu(text), []string{"девопс"}), say("Девопсы не нужны")),
		when(has(brain.normalizeRu(text), []string{"трансформациями"}), replace(text, "трансформациями", "оргиями гомогеев")),
		when(has(brain.normalizeRu(text), []string{"трансформация"}), replace(text, "трансформация", "оргия гомогеев")),
		when(has(brain.normalizeRu(text), []string{"трансформацию"}), replace(text, "трансформацию", "оргию гомогеев")),
		when(has(brain.normalizeRu(text), []string{"трансформации"}), replace(text, "трансформации", "оргии гомогеев")),
		when(has(text, []string{"java", "джаба", "джава"}), say("джава-хуява, а я работаю на го")),
		when(has(text, []string{"блокир"}), say("пусть себе анус заблокируют")),
		random(10, when(has(text, []string{"опять"}), say("не опять, а снова"))),
		when(has(text, []string{"проблем"}), say("у меня есть 5-10 солюшенов этой проблемы")),
	} {
		has, message := opinion.Express(text)
		if has && message != "" {
			return true, message
		}
	}
	return false, ""
}

func (brain *Brain) normalizeRu(text string) string {
	result := text
	for k, v := range brain.GetNormalizationMap() {
		result = strings.Replace(result, k, v, -1)
	}
	return result
}

func (brain *Brain) normalizeEn(text string) string {
	result := text
	for k, v := range brain.GetNormalizationMap() {
		result = strings.Replace(result, v, k, -1)
	}
	return result
}

func (brain *Brain) GetNormalizationMap() map[string]string {
	return brain.memory.normalisationMap
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
	condition bool
	opinion   Opinion
}

func (this ConditionOpinion) Express(text string) (has bool, response string) {
	if this.condition {
		return this.opinion.Express(text)
	} else {
		return false, ""
	}
}

func when(condition bool, opinion Opinion) *ConditionOpinion {
	return &ConditionOpinion{condition: condition, opinion: opinion}
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

func has(text string, values []string) bool {
	for _, value := range values {
		if strings.Contains(text, value) {
			return true
		}
	}
	for _, value := range values {
		if strings.Contains(text, value) {
			return true
		}
	}
	return false
}

func isOf(text string, values []string) bool {
	for _, value := range values {
		if text == value {
			return true
		}
	}
	return false
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
