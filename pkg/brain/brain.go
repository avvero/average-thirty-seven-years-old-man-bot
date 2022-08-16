package brain

import (
	"strconv"
	"strings"

	"github.com/avvero/the_gamers_guild_bot/internal/knowledge"
	"github.com/avvero/the_gamers_guild_bot/internal/utils"
)

type Brain struct {
	memory *Memory
}

func NewBrain() *Brain {
	return &Brain{memory: NewMemory()}
}

func (b *Brain) Decision(chatId int64, text string, rnd bool) (respond bool, response string) {
	if rnd {
		if utils.RandomUpTo(100) == 0 {
			phrase := b.GetSenslessPhrases()[utils.RandomUpTo(len(b.GetSenslessPhrases()))]
			return true, phrase
		}
		if len(text) > 14 && utils.RandomUpTo(100) == 0 {
			phrase := b.khaleesify(text)
			return true, phrase
		}
		if !utils.Contains([]string{"0", "-1001733786877", "245851441", "-578279468"}, strconv.FormatInt(chatId, 10)) {
			return true, "Mr Moony presents his compliments to Professor Snape, and begs him to keep his abnormally large nose out of other people’s business."
		}
	}
	//
	text = strings.ToLower(text)
	if text == "gg" {
		return true, "gg"
	}
	if b.normalizeRu(text) == "нет" {
		return true, "пидора ответ"
	}
	if strings.Contains(text, "morrowind") ||
		strings.Contains(text, "морровинд") ||
		strings.Contains(text, "моровинд") {
		return true, "Morrowind - одна из лучших игр эва"
	}
	if text == "er" ||
		text == "ер" ||
		text == "эр" ||
		strings.Contains(text, "elden ring") ||
		strings.Contains(text, " er ") ||
		strings.Contains(text, " ер ") ||
		strings.Contains(text, " эр ") {
		return true, "Elden Ring - это величие"
	}
	if strings.Contains(text, "купил") {
		return true, "А не пиздишь? Аренда это не покупка"
	}
	if strings.Contains(text, "spotify") || strings.Contains(text, "спотифай") {
		return true, "Эти пидоры Антону косарик должны за подписку"
	}
	if strings.Contains(b.normalizeEn(text), "devops") ||
		strings.Contains(b.normalizeRu(text), "девопс") {
		return true, "Девопсы не нужны"
	}
	if text == "трансформация" ||
		text == "трансформацию" ||
		text == "трансформации" ||
		strings.Contains(b.normalizeRu(text), "трансформация ") ||
		strings.Contains(b.normalizeRu(text), "трансформацию ") ||
		strings.Contains(b.normalizeRu(text), "трансформации ") {
		tokens := map[string]string{
			"трансформация": "оргия гомогеев",
			"трансформацию": "оргию гомогеев",
			"трансформации": "оргии гомогеев",
		}
		result := text
		for k, v := range tokens {
			result = strings.Replace(result, k, v, -1)
		}
		return true, result
	}
	if strings.Contains(text, "java") ||
		strings.Contains(text, "джаба") ||
		strings.Contains(text, "джава") {
		return true, "джава-хуява, а я работаю на го"
	}
	return false, ""
}

func (b *Brain) normalizeRu(text string) string {
	result := text
	for k, v := range b.GetNormalizationMap() {
		result = strings.Replace(result, k, v, -1)
	}
	return result
}

func (b *Brain) normalizeEn(text string) string {
	result := text
	for k, v := range b.GetNormalizationMap() {
		result = strings.Replace(result, v, k, -1)
	}
	return result
}

func (b *Brain) khaleesify(text string) string {
	result := strings.ToLower(text)
	for _, k := range b.GetMokingMapKeys() {
		result = strings.Replace(result, k, b.GetMockingMap()[k], -1)
	}
	return result
}

func (b *Brain) GetSenslessPhrases() []string {
	return b.memory.senselessPhrases
}

func (b *Brain) GetMockingMap() map[string]string {
	return b.memory.mockingMap
}

func (b *Brain) GetMokingMapKeys() []string {
	return b.memory.mockingMapKeys
}

func (b *Brain) GetNormalizationMap() map[string]string {
	return b.memory.normalisationMap
}

func (b *Brain) RememberAll() *Brain {
	b.memory.SetSenslessPhrases(knowledge.SenselessPhrases)
	b.memory.SetMockingMap(knowledge.MockingMap)
	b.memory.SetNormalizationMap(knowledge.NormalisationMap)
	return b
}
